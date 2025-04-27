package gatewayImpl

import (
	"fmt"
	"log"

	"example.com/infrahandson/internal/interface/gateway"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
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

func (i *SQLiteInitializerImpl) InitSchema(db *sqlx.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	public_id TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT
);
CREATE TABLE IF NOT EXISTS rooms (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  public_id TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS room_members (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  room_id INTEGER NOT NULL,
  user_id TEXT NOT NULL,
  FOREIGN KEY (room_id) REFERENCES rooms(id)
);`

	_, err := db.Exec(schema)
	return err
}
