package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	shortlinksTable := `
	CREATE TABLE IF NOT EXISTS shortlinks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code VARCHAR(6) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
		is_disabled BOOLEAN NOT NULL DEFAULT FALSE,
		total_visits INTEGER NOT NULL DEFAULT 0,
		total_duration INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_shortlinks_code ON shortlinks(code);
	`

	visitsTable := `
	CREATE TABLE IF NOT EXISTS visits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		shortlink_id INTEGER NOT NULL,
		ip_address VARCHAR(45) NOT NULL,
		country VARCHAR(64),
		city VARCHAR(128),
		latitude DECIMAL(10,6),
		longitude DECIMAL(10,6),
		geo_precision VARCHAR(16),
		geo_status VARCHAR(16),
		user_agent TEXT,
		os_type VARCHAR(32),
		os_version VARCHAR(32),
		browser_type VARCHAR(32),
		browser_version VARCHAR(32),
		device_type VARCHAR(16),
		visit_duration INTEGER,
		visit_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		exit_time DATETIME,
		referer TEXT,
		FOREIGN KEY (shortlink_id) REFERENCES shortlinks(id)
	);
	CREATE INDEX IF NOT EXISTS idx_visits_shortlink_id ON visits(shortlink_id);
	CREATE INDEX IF NOT EXISTS idx_visits_visit_time ON visits(visit_time);
	`

	settingsTable := `
	CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		[key] VARCHAR(64) UNIQUE NOT NULL,
		value TEXT,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	tables := []string{shortlinksTable, visitsTable, settingsTable}
	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}
