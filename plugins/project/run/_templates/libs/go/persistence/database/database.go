package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	maxBatch  = 50
	inMemFile = "file::memory:?cache=private"
)

func ReadyStateDB(databaseURL string, models ...interface{}) (*gorm.DB, error) {
	db, err := newPostgresGorm(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("getting postgres connection: %w", err)
	}
	err = MigrateDbSchema(db, models...)
	if err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}
	return db, nil
}

func newPostgresGorm(databaseUrl string) (*gorm.DB, error) {
	return gorm.Open(postgres.New(postgres.Config{
		DSN: databaseUrl,
	}), &gorm.Config{CreateBatchSize: maxBatch})
}

func ReadyTestDB(models ...interface{}) (*gorm.DB, error) {
	db, err := inMemGorm()
	if err != nil {
		return nil, fmt.Errorf("creating in memory database: %w", err)
	}
	err = MigrateDbSchema(db, models...)
	if err != nil {
		return nil, fmt.Errorf("migrating database schema: %w", err)
	}
	return db, nil
}

func inMemGorm() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(inMemFile), &gorm.Config{CreateBatchSize: maxBatch})
}

func NewPostgres(databaseUrl string) (*sql.DB, error) {
	return sql.Open("postgres", databaseUrl)
}

func MigrateDbSchema(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	return nil
}
