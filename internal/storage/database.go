package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dataDir string) (*Database, error) {
	dbPath := filepath.Join(dataDir, "rag-therapist.db")
	
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return database, nil
}

func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS documents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_name TEXT NOT NULL,
		file_path TEXT NOT NULL UNIQUE,
		file_size INTEGER NOT NULL,
		content_hash TEXT NOT NULL,
		uploaded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		processed_at DATETIME,
		status TEXT NOT NULL DEFAULT 'pending'
	);

	CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(status);
	CREATE INDEX IF NOT EXISTS idx_documents_uploaded_at ON documents(uploaded_at);
	CREATE INDEX IF NOT EXISTS idx_documents_content_hash ON documents(content_hash);
	`

	_, err := d.db.Exec(query)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}