package postgres

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func InitDB(databaseURL string, log *slog.Logger) (*sql.DB, error) {
	var err error

	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Error("failed to open DB connection", "error", err)
		return conn, err
	}

	if err := conn.Ping(); err != nil {
		log.Error("failed to ping DB", "error", err)
		_ = conn.Close()
		return nil, err
	}

	log.Info("database connected!")
	return conn, nil
}

func CloseDB(db *sql.DB, log *slog.Logger) {
	if err := db.Close(); err != nil {
		log.Error("error closing DB", "error", err)
	}
	log.Info("database closed!")
}
