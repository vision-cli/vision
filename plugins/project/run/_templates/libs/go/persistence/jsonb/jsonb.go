package jsonb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DB[M Model] struct {
	db      DBTX
	table   string
	archive string
	hard    bool
	conds   conditions
}

type Model interface {
	ID() string
}

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func newDB[M Model](db DBTX, table string) (*DB[M], error) {
	var err error
	ctx := context.Background()
	jsb := &DB[M]{
		db:      db,
		table:   table,
		archive: fmt.Sprintf("_%s", table),
		hard:    false,
	}

	if err = jsb.createTable(ctx); err != nil {
		return nil, fmt.Errorf("creating database schema: %w", err)
	}

	if err = jsb.indexTable(ctx); err != nil {
		return nil, fmt.Errorf("indexing jsonb table: %w", err)
	}

	return jsb, err
}

func (d *DB[M]) createTable(ctx context.Context) error {
	stmts := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %[1]s (id VARCHAR(255) PRIMARY KEY,doc JSONB,deleted_at TIMESTAMP);
CREATE OR REPLACE VIEW %[2]s AS SELECT id,doc FROM %[1]s WHERE deleted_at IS NULL;
CREATE OR REPLACE RULE _soft_delete AS ON DELETE TO %[2]s DO INSTEAD (UPDATE %[1]s SET deleted_at=now() WHERE id=old.id);`,
		d.archive, d.table)

	_, err := d.db.ExecContext(ctx, stmts)
	if err != nil {
		return fmt.Errorf("creating table: %w", err)
	}

	return nil
}

func (d *DB[M]) indexTable(ctx context.Context) error {
	indexStmt := fmt.Sprintf(
		"CREATE INDEX IF NOT EXISTS %s_gin ON %s USING GIN (doc jsonb_path_ops)",
		d.table, d.archive)

	_, err := d.db.ExecContext(ctx, indexStmt)
	if err != nil {
		return fmt.Errorf("creating index for table %q: %w", d.table, err)
	}
	return nil
}

func (d *DB[M]) dropTable(ctx context.Context) error {
	dropStmt := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", d.archive)

	_, err := d.db.ExecContext(ctx, dropStmt)
	if err != nil {
		return fmt.Errorf("dropping table %q: %w", d.table, err)
	}
	return nil
}

func (d *DB[M]) resetTable(ctx context.Context) error {
	resetStmt := fmt.Sprintf("TRUNCATE TABLE %s", d.archive)

	_, err := d.db.ExecContext(ctx, resetStmt)
	if err != nil {
		return fmt.Errorf("resetting table %q: %w", d.table, err)
	}
	return nil
}

// CreateInBatches inserts models into the database in batches of batchSize
func (d *DB[M]) CreateInBatches(ctx context.Context, models []M, batchSize int) error {
	return d.ensureTx(func(tx *DB[M]) error {
		return batch(models, batchSize, func(batch []M) error {
			var err error
			values := make([]string, 0, len(batch))

			for _, model := range batch {
				var row Row
				row, err = parse(model.ID(), model)
				if err != nil {
					return err
				}

				values = append(values, fmt.Sprintf("('%s','%s')", row.id, row.doc))
			}

			stmt := fmt.Sprintf("INSERT INTO %s (id,doc) VALUES %s", tx.table, strings.Join(values, ","))

			_, err = tx.db.ExecContext(ctx, stmt)
			if err != nil {
				return fmt.Errorf("inserting batch: %w", err)
			}
			return nil
		})
	})
}

// SaveInBatches updates the existing models in the database in batches of batchSize.
// Any models not present in the database will be ignored
func (d *DB[M]) SaveInBatches(ctx context.Context, models []M, batchSize int) error {
	return d.ensureTx(func(tx *DB[M]) error {
		return batch(models, batchSize, func(batch []M) error {
			var err error
			values := make([]string, 0, len(batch))

			for _, model := range batch {
				var row Row
				row, err = parse(model.ID(), model)
				if err != nil {
					return err
				}

				values = append(values, fmt.Sprintf("('%s','%s'::jsonb)", row.id, row.doc))
			}

			stmt := fmt.Sprintf("UPDATE %s AS t SET doc=u.doc FROM (VALUES %s) AS u(id,doc) WHERE u.id=t.id",
				tx.table, strings.Join(values, ","))

			_, err = tx.db.ExecContext(ctx, stmt)
			if err != nil {
				return fmt.Errorf("updating batch: %w", err)
			}
			return nil
		})
	})
}

// DeleteInBatches removes all entries corresponding to ids in batches of batchSize.
func (d *DB[M]) DeleteInBatches(ctx context.Context, ids []string, batchSize int) error {
	return d.ensureTx(func(tx *DB[M]) error {
		return batch(ids, batchSize, func(batch []string) error {
			stmt := fmt.Sprintf("DELETE FROM %s WHERE id IN ('%s')", tx.target(), strings.Join(ids, "','"))

			_, err := tx.db.ExecContext(ctx, stmt)
			if err != nil {
				return fmt.Errorf("updating batch: %w", err)
			}
			return nil
		})
	})
}

// batch executes fn for each slice of max length batchSize in turn for all input values.
// batch continues until all of input has been computed or an error occurs.
func batch[T any](input []T, batchSize int, fn func(batch []T) error) error {
	var err error
	l := len(input)

	for i := 0; i < l; i += batchSize {
		end := i + batchSize
		if end > l {
			end = l
		}

		err = fn(input[i:end])
		if err != nil {
			return fmt.Errorf("executing batch %d: %w", i/batchSize, err)
		}
	}

	return nil
}

// Create inserts a new entry for model in the database.
func (d *DB[M]) Create(ctx context.Context, model M) error {
	return d.CreateAny(ctx, model.ID(), model)
}

// CreateAny inserts a new entry for any given structure (doc) in the database.
// Prefer Create for general use. Maps and slices can be used as filters and control structures
// to help look up and extrapolate data within models, approximating subscripting.
func (d *DB[M]) CreateAny(ctx context.Context, id string, doc any) error {
	row, err := parse(id, doc)
	if err != nil {
		return fmt.Errorf("parsing %T into json: %w", doc, err)
	}

	insert := fmt.Sprintf("INSERT INTO %s (id,doc) VALUES ('%s','%s')", d.table, row.id, row.doc)

	_, err = d.db.ExecContext(ctx, insert)
	if err != nil {
		return fmt.Errorf("executing statement %q: %w", insert, err)
	}

	return nil
}

// Find retrieves a model with the matching id from the database.
func (d *DB[M]) Find(ctx context.Context, id string) (M, error) {
	var model M
	err := d.FindAny(ctx, id, &model)

	return model, err
}

// FindAny retrieves any existing structure from the database with id, populating doc.
// Prefer Find for general use.
func (d *DB[M]) FindAny(ctx context.Context, id string, doc any) error {
	var row Row
	selectStmt := fmt.Sprintf("SELECT id,doc FROM %s WHERE id='%s'", d.table, id)

	err := d.db.QueryRowContext(ctx, selectStmt).Scan(&row.id, &row.doc)
	if err != nil {
		return fmt.Errorf("scanning result from query %q into Row struct: %w", selectStmt, err)
	}

	return row.decode(doc)
}

// Save updates an existing entry for model in the database.
// Commonly used after Find when subsequent changes have been made to the retrieved model.
func (d *DB[M]) Save(ctx context.Context, model M) error {
	return d.SaveAny(ctx, model.ID(), model)
}

// Save updates any existing structure in the database by id, replacing the entry with doc.
// Prefer Save for general use.
func (d *DB[M]) SaveAny(ctx context.Context, id string, doc any) error {
	row, err := parse(id, doc)
	if err != nil {
		return fmt.Errorf("parsing %T into json: %w", doc, err)
	}

	updateStmt := fmt.Sprintf("UPDATE %s SET doc='%s' WHERE id='%s'", d.table, row.doc, row.id)

	res, err := d.db.ExecContext(ctx, updateStmt)
	if err != nil {
		return fmt.Errorf("executing statement %q: %w", updateStmt, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("retrieving affected rows: %w", err)
	}
	if affected != 1 {
		return fmt.Errorf("unable to update entry with id %q", id)
	}

	return nil
}

// Delete removes entry with id from the database.
func (d *DB[M]) Delete(ctx context.Context, id string) error {
	deleteStmt := fmt.Sprintf("DELETE FROM %s WHERE id='%s'", d.target(), id)

	_, err := d.db.ExecContext(ctx, deleteStmt)
	if err != nil {
		return fmt.Errorf("executing statement %q: %w", deleteStmt, err)
	}

	return nil
}

// Hard returns a DB that hard deletes records. By default DB soft deletes records.
func (d DB[M]) Hard() *DB[M] {
	d.hard = true
	return &d
}

// return target on which delete operations should be performed
func (d *DB[M]) target() string {
	if d.hard {
		return d.archive
	}
	return d.table
}

// Models returns all models from the database where any optional conditions are met.
func (d *DB[M]) Models(ctx context.Context) ([]M, error) {
	c, err := d.conds.clause()
	if err != nil {
		return nil, err
	}

	models, err := d.models(ctx, selectStmt(d.table, c))
	if err != nil {
		return nil, fmt.Errorf("finding models from rows: %w", err)
	}

	return models, nil
}

// Count returns the number of rows from the database where any optional conditions are met
func (d *DB[M]) Count(ctx context.Context) (int, error) {
	var count int
	c, err := d.conds.clause()
	if err != nil {
		return 0, err
	}

	countStmt := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", d.table, c)

	err = d.db.QueryRowContext(ctx, countStmt).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting rows: %w", err)
	}

	return count, err
}

// Paginate returns a number of models from the database wuth any given conditions up to size limit.
// The queryFunc constructs the query based on various parameters e.g. offset or keysey pegination.
func (d *DB[M]) Paginate(ctx context.Context, queryFunc pageQuery) ([]M, error) {
	pageStmt, err := queryFunc(d.table, d.conds)
	if err != nil {
		return nil, err
	}

	models, err := d.models(ctx, pageStmt)
	if err != nil {
		return nil, fmt.Errorf("finding models from rows: %w", err)
	}

	return models, nil
}

type pageQuery func(string, conditions) (string, error)

func PageNum(n int, limit int) pageQuery {
	offset := n * limit
	return func(table string, conds conditions) (string, error) {
		c, err := conds.clause()
		if err != nil {
			return "", err
		}

		offsetClause := ""
		if n > 0 {
			offsetClause = fmt.Sprintf(" OFFSET %d", offset)
		}
		return fmt.Sprintf("%s ORDER BY id ASC LIMIT %d%s",
			selectStmt(table, c), limit, offsetClause), nil
	}
}

func PageAfter(id string, limit int) pageQuery {
	return func(table string, conds conditions) (string, error) {
		conds.add(fmt.Sprintf("id > '%s'", id))
		c, err := conds.clause()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s ORDER BY id ASC LIMIT %d",
			selectStmt(table, c), limit), nil
	}
}

func PageBefore(id string, limit int) pageQuery {
	return func(table string, conds conditions) (string, error) {
		conds.add(fmt.Sprintf("id < '%s'", id))
		c, err := conds.clause()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("SELECT * FROM (%s ORDER BY id DESC LIMIT %d) AS prev_page ORDER BY id ASC",
			selectStmt(table, c), limit), nil
	}
}

func selectStmt(table, whereClause string) string {
	return fmt.Sprintf("SELECT * FROM %s%s", table, whereClause)
}

func (d *DB[M]) models(ctx context.Context, stmt string) ([]M, error) {
	sqlRows, err := d.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("executing statement %q: %w", stmt, err)
	}
	defer sqlRows.Close()

	models := make([]M, 0)
	for sqlRows.Next() {
		var row Row
		if err = sqlRows.Scan(&row.id, &row.doc); err != nil {
			return nil, fmt.Errorf("scanning sql rows: %w", err)
		}

		var m M
		if err = row.decode(&m); err != nil {
			return nil, fmt.Errorf("decoding row (%+v): %w", row, err)
		}

		if m.ID() == row.id {
			models = append(models, m)
		}
	}

	return models, nil
}

// Begin opens a transaction providing a transaction is not already in progress.
func (d DB[M]) Begin() (*DB[M], error) {
	db, ok := d.db.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("cannot begin transaction in existing transaction")
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}

	d.db = tx
	return &d, nil
}

// Rollback aborts the current transaction
func (d *DB[M]) Rollback() error {
	tx, ok := d.db.(*sql.Tx)
	if !ok {
		return fmt.Errorf("cannot rollback, transaction not in progress")
	}
	return tx.Rollback()
}

// Commit commits the current transaction
func (d *DB[M]) Commit() error {
	tx, ok := d.db.(*sql.Tx)
	if !ok {
		return fmt.Errorf("cannot commit, transaction not in progress")
	}
	return tx.Commit()
}

// Transaction starts a transaction as a block.
// Any returned errors will abort the transaction.
// If no errors are returned, the transaction is then committed.
func (d *DB[M]) Transaction(fn func(tx *DB[M]) error) error {
	var err error
	panicked := true

	tx, err := d.Begin()
	if err != nil {
		return fmt.Errorf("preparing transaction: %w", err)
	}

	defer func() {
		// Make sure to rollback on panic, or error
		if panicked || err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("%v: %w", rbErr, err)
			}
		}
	}()

	if err = fn(tx); err == nil {
		panicked = false
		return tx.Commit()
	}

	panicked = false
	return err
}

// ensureTx wraps fn in a transaction if a transaction is not already in progress.
// Allows a method to be used singularly as well as in series as part of a larger method.
func (d *DB[M]) ensureTx(fn func(tx *DB[M]) error) error {
	_, ok := d.db.(*sql.Tx)
	if !ok {
		return d.Transaction(fn)
	}
	return fn(d)
}
