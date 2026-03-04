package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) {
	var err error
	// 启用外键支持（可选）
	DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		log.Fatal("Failed to connect to SQLite:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping SQLite:", err)
	}

	// 启用 WAL 模式提升并发性能
	DB.Exec("PRAGMA journal_mode=WAL;")

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS monitored_sites (
            id TEXT PRIMARY KEY,
            url TEXT NOT NULL,
            name TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            is_active INTEGER DEFAULT 1
        );`,
		`CREATE TABLE IF NOT EXISTS change_records (
            id TEXT PRIMARY KEY,
            site_id TEXT NOT NULL,
            change_type TEXT NOT NULL,
            file_path TEXT,
            change_diff TEXT,
            snapshot_hash TEXT,
            detected_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE TABLE IF NOT EXISTS file_contents (
            id TEXT PRIMARY KEY,
            site_id TEXT NOT NULL,
            file_path TEXT,
            content_hash TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            change_record_id TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS hook_scripts (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            script_type TEXT NOT NULL,
            script_content TEXT NOT NULL,
            version TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE TABLE IF NOT EXISTS site_hooks (
            id TEXT PRIMARY KEY,
            site_id TEXT NOT NULL,
            script_id TEXT NOT NULL,
            enabled INTEGER DEFAULT 1,
            config TEXT,
            priority INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
	log.Println("SQLite tables initialized.")
}
