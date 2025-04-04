package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	// This IS being used!
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

/**
* Sets up the Database connection and provides its access as a singleton to
* the entire application.
**/
func InitDB() *sqlx.DB {
	// construct the db connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("constructed dsn", dsn)

	// pass the db connection string to connect to our database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	fmt.Printf("\nConnected to the database successfully.\n\n")

	// set global instance for the database
	DB = db
	return DB
}
