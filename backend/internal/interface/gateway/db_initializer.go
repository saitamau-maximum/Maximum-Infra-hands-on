package gateway

import "github.com/jmoiron/sqlx"

type DBInitializer interface {
	Init() (*sqlx.DB, error)
	InitSchema(db *sqlx.DB) error
}
