package mysqlgatewayimpl

import (
	"errors"

	"example.com/infrahandson/internal/interface/gateway"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLInitializer struct {
	dsn            string
	migrationsPath string
}

type NewMySQLInitializerParams struct {
	DSN            *string
	MigrationsPath *string
}

func (p *NewMySQLInitializerParams) Validate() error {
	if p.DSN == nil {
		return errors.New("DSN is required")
	}
	if p.MigrationsPath == nil {
		return errors.New("migrationsPath is required")
	}
	return nil
}

func NewMySQLInitializer(params *NewMySQLInitializerParams) gateway.DBInitializer {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	// Validateに通るなら，*params.DSNや*parms.MigrationPathがnil pointerにならないことが保証される
	return &MySQLInitializer{
		dsn:            *params.DSN,
		migrationsPath: *params.MigrationsPath,
	}
}

func (i *MySQLInitializer) Init() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", i.dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (i *MySQLInitializer) InitSchema(db *sqlx.DB) error {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return err
	}

	path := "file://" + i.migrationsPath

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
