package gateway_impl

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"example.com/webrtc-practice/internal/interface/gateway"
)

type SQLiteInitializerImpl struct {
	Path string // e.g., ":memory:" or "data/test.db"
}

func NewSQLiteInitializer(path string) gateway.DBInitializer {
	return &SQLiteInitializerImpl{Path: path}
}

func (i *SQLiteInitializerImpl) Init() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", i.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite: %w", err)
	}

	log.Printf("SQLite connected at %s\n", i.Path)
	return db, nil
}
