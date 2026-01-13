package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(connStr string) error {
	var err error

	// Retry loop for DB connection (wait for docker container)
	for i := 0; i < 5; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}

		if err = DB.Ping(); err == nil {
			break
		}

		log.Printf("Failed to connect to DB, retrying in 2s... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}

	log.Println("Connected to Database")
	return Migrate()
}

func Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS chats (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
			role VARCHAR(20) NOT NULL, 
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("migration failed for query: %s \n error: %v", query, err)
		}
	}
	return nil
}
