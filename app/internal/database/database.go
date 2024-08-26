package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	GetDB() *sql.DB
	Close() error
}

type SQLiteDatabase struct {
	db *sql.DB
}

func New(dbPath string) (Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	database := &SQLiteDatabase{db: db}

	if err := database.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return database, nil
}

func (d *SQLiteDatabase) GetDB() *sql.DB {
	return d.db
}

func (d *SQLiteDatabase) Close() error {
	return d.db.Close()
}

func (d *SQLiteDatabase) runMigrations() error {
	driver, err := sqlite3.WithInstance(d.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	return nil
}
