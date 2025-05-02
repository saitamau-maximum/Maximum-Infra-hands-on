package sqlitegatewayimpl

import (
	"errors"
	"log"

	"example.com/infrahandson/internal/interface/gateway"
	"github.com/golang-migrate/migrate/v4"
	sqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type SQLiteInitializerImpl struct {
	path           string // e.g., ":memory:" or "data/test.db"
	migrationsPath string // e.g., "file://migrations"
}

type NewSQLiteInitializerParams struct {
	Path           string // e.g., ":memory:" or "data/test.db"
	MigrationsPath string // e.g., "file://migrations"
}

func (p *NewSQLiteInitializerParams) Validate() error {
	if p.Path == "" {
		return errors.New("path is required")
	}
	if p.MigrationsPath == "" {
		return errors.New("migrationsPath is required")
	}
	return nil
}

func NewSQLiteInitializer(p *NewSQLiteInitializerParams) gateway.DBInitializer {
	if err := p.Validate(); err != nil {
		panic(err)
	}
	return &SQLiteInitializerImpl{
		path:           p.Path,
		migrationsPath: p.MigrationsPath,
	}
}

func (i *SQLiteInitializerImpl) Init() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", i.path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("SQLite connected at %s\n", i.path)
	return db, nil
}

func (i *SQLiteInitializerImpl) InitSchema(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	path := "file://" + i.migrationsPath

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	// Upマイグレーション実行
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
