package repository

import (
	"database/sql"
	"log"
)

func InitDB(connPath string) *sql.DB {
	var err error

	conn, err := sql.Open("postgres", connPath)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Database connected!")
	return conn
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}
	log.Println("Database closed!")
}
