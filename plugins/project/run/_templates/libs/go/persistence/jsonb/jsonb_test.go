package jsonb

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type jsonbSuite struct {
	suite.Suite
	db   *DB[*Mod]
	data []*Mod
}

type Mod struct {
	Id   int
	Info string
	Tags []string
}

func (m *Mod) ID() string {
	return strconv.Itoa(m.Id)
}

func TestJsonbSuite(t *testing.T) {
	suite.Run(t, new(jsonbSuite))
}

func (j *jsonbSuite) SetupSuite() {
	db, _, err := NewTest[*Mod](dbUrl)
	require.NoError(j.T(), err)

	j.db = db
	j.data = []*Mod{
		{0, "mod1", nil},
		{1, "mod2", []string{"test"}},
		{2, "mod3", []string{"test", "foo"}},
	}
}

func (j *jsonbSuite) TearDownSuite() {
	assert.NoError(j.T(), j.db.dropTable(context.Background()))
}

// All tests depend on insert. Ensure TestCreate passes first.
func (j *jsonbSuite) SetupDB(models []*Mod) (context.Context, *DB[*Mod]) {
	ctx := context.Background()
	tx, err := j.db.Begin()
	require.NoError(j.T(), err)

	for _, m := range models {
		err = tx.Create(ctx, m)
		require.NoError(j.T(), err)
	}

	return ctx, tx
}

func (j *jsonbSuite) TestCreate() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	j.assertCount(ctx, tx, tx.table, 3)

	// cannot insert again
	err = tx.Create(ctx, j.data[2])
	require.Error(j.T(), err)
}

func (j *jsonbSuite) TestCreateInBatches() {
	var err error
	ctx, tx := j.SetupDB(nil)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	err = tx.CreateInBatches(ctx, j.data, 2)
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, len(j.data))
}

func (j *jsonbSuite) TestFind() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	// models are correctly returned intact
	for _, mod := range j.data {
		actual, err := tx.Find(ctx, mod.ID())
		require.NoError(j.T(), err)
		assert.Equal(j.T(), *mod, *actual)
	}
}

func (j *jsonbSuite) TestModels() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	// should not be returned
	err = tx.CreateAny(ctx, "map", map[int]string{1: "test"})
	require.NoError(j.T(), err)

	models, err := tx.Models(ctx)
	require.NoError(j.T(), err)
	assert.Equal(j.T(), len(j.data), len(models))

	// models appear once each
	for _, d := range j.data {
		var count int

		for _, model := range models {
			if model.ID() == d.ID() {
				count++
				assert.Equal(j.T(), *d, *model)
			}
		}
		assert.Equal(j.T(), 1, count)
	}
}

func (j *jsonbSuite) TestCount() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	count, err := tx.Count(ctx)
	require.NoError(j.T(), err)
	assert.Equal(j.T(), len(j.data), count)
}

// depends on Find()
func (j *jsonbSuite) TestSave() {
	var err error
	ctx, tx := j.SetupDB([]*Mod{
		j.data[0],
		j.data[1],
	})

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	updated := *j.data[1]
	updated.Info = "updated"

	err = tx.Save(ctx, &updated)
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 2)

	actual, err := tx.Find(ctx, j.data[1].ID())
	require.NoError(j.T(), err)
	assert.Equal(j.T(), updated, *actual)

	// cannot save if not present in db
	err = tx.Save(ctx, j.data[2])
	require.Error(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 2)
}

// depends on Find()
func (j *jsonbSuite) TestSaveInBatches() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	updated := make([]*Mod, len(j.data))
	copy(updated, j.data)
	for i, m := range updated {
		update := *m
		updated[i] = &update
	}
	// skip middle, see if unaffected
	updated[0].Info = "update0"
	updated[2].Info = "update2"

	err = tx.SaveInBatches(ctx, []*Mod{updated[0], updated[2]}, 3)
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 3)

	var actual *Mod
	for _, u := range updated {
		actual, err = tx.Find(ctx, u.ID())
		require.NoError(j.T(), err)
		assert.Equal(j.T(), *u, *actual)
	}

	// skips save if not present in db
	err = tx.SaveInBatches(ctx, []*Mod{{Id: 3, Info: "", Tags: nil}}, 3)
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 3)
}

// depends on Find()
func (j *jsonbSuite) TestDelete() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	j.assertCount(ctx, tx, tx.table, 3)

	// Idempotency: deleting a non-existent entry.
	// User will believe entry now doesn't exist
	// which is always true regardless of if it existed before.
	for i := 0; i < 2; i++ {
		err = tx.Delete(ctx, j.data[0].ID())
		require.NoError(j.T(), err)
	}

	j.assertCount(ctx, tx, tx.table, 2)
	j.assertCount(ctx, tx, tx.archive, 3)

	actual, err := tx.Find(ctx, j.data[1].ID())
	require.NoError(j.T(), err)
	assert.Equal(j.T(), *j.data[1], *actual)

	for i := 0; i < 2; i++ {
		err = tx.Hard().Delete(ctx, j.data[0].ID())
		require.NoError(j.T(), err)
	}

	j.assertCount(ctx, tx, tx.table, 2)
	j.assertCount(ctx, tx, tx.archive, 2)

	err = tx.Hard().Delete(ctx, j.data[1].ID())
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 1)
	j.assertCount(ctx, tx, tx.archive, 1)
}

// depends on Find()
func (j *jsonbSuite) TestDeleteInBatches() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	// ignores ids that don't exist
	err = tx.DeleteInBatches(ctx, []string{j.data[0].ID(), "999", j.data[2].ID()}, 2)
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 1)
	j.assertCount(ctx, tx, tx.archive, 3)

	actual, err := tx.Find(ctx, j.data[1].ID())
	require.NoError(j.T(), err)
	assert.Equal(j.T(), *j.data[1], *actual)

	err = tx.Hard().DeleteInBatches(ctx, []string{j.data[0].ID(), "999", j.data[2].ID()}, 2)
	require.NoError(j.T(), err)

	j.assertCount(ctx, tx, tx.table, 1)
	j.assertCount(ctx, tx, tx.archive, 1)
}

func (j *jsonbSuite) TestMatch() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	for _, test := range []struct {
		in map[string]any
		ex int
	}{
		// specific field
		{in: map[string]any{"Info": "mod2"}, ex: 1},
		// all entries
		{in: map[string]any{}, ex: 3},
		// any tags
		{in: map[string]any{"Tags": []string{}}, ex: 2},
		// no tags
		{in: map[string]any{"Tags": nil}, ex: 1},
		// includes test tag
		{in: map[string]any{"Tags": []string{"test"}}, ex: 2},
		// includes both tags
		{in: map[string]any{"Tags": []string{"test", "foo"}}, ex: 1},
	} {
		var models []*Mod
		models, err = tx.Match(test.in).Models(ctx)
		require.NoError(j.T(), err)
		assert.Equal(j.T(), test.ex, len(models))

		var count int
		count, err = tx.Match(test.in).Count(ctx)
		require.NoError(j.T(), err)
		assert.Equal(j.T(), test.ex, count)
	}
}

func (j *jsonbSuite) TestFilter() {
	var err error
	ctx, tx := j.SetupDB(j.data)

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	for _, test := range []struct {
		in map[string][]string
		ex int
	}{
		// specific field
		{in: map[string][]string{"Info": {"mod2"}}, ex: 1},
		// all entries
		{in: map[string][]string{"Info": nil, "Tags": nil}, ex: 3},
		// either of two
		{in: map[string][]string{"Info": {"mod1", "mod2"}}, ex: 2},
		// intersection of filters
		{in: map[string][]string{"Id": {"1", "2"}, "Info": {"mod1", "mod2"}}, ex: 1},
		// match even whole slices, but must be exact format, not subset
		{in: map[string][]string{"Tags": {`["test"]`, `["test", "foo"]`}}, ex: 2},
	} {
		var models []*Mod
		models, err = tx.Filter(test.in).Models(ctx)
		require.NoError(j.T(), err)
		assert.Equal(j.T(), test.ex, len(models))

		var count int
		count, err = tx.Filter(test.in).Count(ctx)
		require.NoError(j.T(), err)
		assert.Equal(j.T(), test.ex, count)
	}
}

func (j *jsonbSuite) TestSearch() {
	var err error
	ctx, tx := j.SetupDB(append(j.data,
		&Mod{Id: 11, Info: "anotherMod", Tags: []string{"foobar"}},
	))

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	for _, test := range []struct {
		in map[string]string
		ex int
	}{
		// one field
		{in: map[string]string{"Info": "mod1"}, ex: 1},
		// all entries, caps insensitive
		{in: map[string]string{"Info": "MOD"}, ex: 4},
		// wildcard, must be a character after "mod"
		{in: map[string]string{"Info": "mod_"}, ex: 3},
		// union of searches
		{in: map[string]string{"Id": "1", "Tags": "FOO"}, ex: 3},
		// blank field are ignored
		{in: map[string]string{"Id": "", "Info": "another"}, ex: 1},
	} {
		var models []*Mod
		models, err = tx.Search(test.in).Models(ctx)
		require.NoError(j.T(), err)
		assert.Equal(j.T(), test.ex, len(models))

		var count int
		count, err = tx.Search(test.in).Count(ctx)
		require.NoError(j.T(), err)
		assert.Equal(j.T(), test.ex, count)
	}
}

func (j *jsonbSuite) TestMultipleConds() {
	var err error

	extras := []*Mod{
		{Id: 11, Info: "anotherMod", Tags: []string{"test", "foobar"}},
		{Id: 12, Info: "another", Tags: []string{"foobar"}},
		{Id: 13, Info: "andMore", Tags: []string{"test", "foobar"}},
	}

	ctx, tx := j.SetupDB(append(j.data, extras...))

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	match := map[string]any{"Tags": []string{"test"}}
	filter := map[string][]string{"Id": {"0", "1", "11", "13"}}
	search := map[string]string{"Info": "mod"}
	expected := []*Mod{j.data[1], extras[0]}

	var models []*Mod
	models, err = tx.Match(match).Filter(filter).Search(search).Models(ctx)
	require.NoError(j.T(), err)
	assert.Equal(j.T(), expected, models)
}

func (j *jsonbSuite) TestPaginate() {
	var err error

	extras := []*Mod{
		{Id: 3, Info: "mod3", Tags: []string{"foo"}},
		{Id: 4, Info: "mod4", Tags: []string{"bar"}},
		{Id: 5, Info: "test", Tags: []string{"test"}},
		{Id: 6, Info: "mod6", Tags: []string{"test"}},
	}

	ctx, tx := j.SetupDB(append(j.data, extras...))

	defer func() {
		err = tx.Rollback()
		require.NoError(j.T(), err)
	}()

	limit := 3

	models, err := tx.Paginate(ctx, PageNum(0, limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), limit, len(models))
	assert.Equal(j.T(), j.data, models)

	models, err = tx.Paginate(ctx, PageNum(1, limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), limit, len(models))
	assert.Equal(j.T(), extras[:3], models)

	after0, err := tx.Paginate(ctx, PageAfter(j.data[2].ID(), limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), limit, len(after0))
	assert.Equal(j.T(), models, after0)

	models, err = tx.Paginate(ctx, PageNum(2, limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), 1, len(models))
	assert.Equal(j.T(), extras[3:], models)

	after1, err := tx.Paginate(ctx, PageAfter(extras[2].ID(), limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), 1, len(after1))
	assert.Equal(j.T(), models, after1)

	before2, err := tx.Paginate(ctx, PageBefore(extras[3].ID(), limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), limit, len(before2))
	assert.Equal(j.T(), after0, before2)

	withConds, err := tx.Search(map[string]string{"Info": "mod"}).Paginate(ctx, PageBefore(extras[3].ID(), limit))
	require.NoError(j.T(), err)
	assert.Equal(j.T(), limit, len(withConds))
	assert.Equal(j.T(), append(j.data[2:], extras[:2]...), withConds)
}

func (j *jsonbSuite) assertCount(ctx context.Context, tx *DB[*Mod], table string, expected int) {
	var actual int
	rows := tx.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(id) FROM %s", table))

	err := rows.Scan(&actual)
	require.NoError(j.T(), err)

	assert.Equal(j.T(), expected, actual)
}
