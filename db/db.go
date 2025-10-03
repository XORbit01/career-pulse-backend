package db

import (
	"database/sql"
	"fmt"

	"github.com/XORbit01/jobseeker-backend/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func Connect(cfg config.DBConfig, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	var connStr string

	// Use full DATABASE_URL if available
	if cfg.DSN != "" {
		connStr = cfg.DSN
	} else {
		// Build from individual DB_* fields
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
		)
	}

	// Connect to the DB
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to open database connection: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("❌ failed to ping database: %w", err)
	}

	// Optional performance tuning
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	return db, nil
}
