package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const dbPath = "store/todos.db"

var db *sql.DB

func Init() error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return err
	}

	dsn := dbPath + "?_pragma=journal_mode(WAL)&_pragma=foreign_keys(on)&_pragma=busy_timeout(5000)"
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}

	db = conn
	return migrate()
}

func migrate() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			uuid       TEXT PRIMARY KEY,
			todo       TEXT NOT NULL,
			completed  INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);
	`)
	return err
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
