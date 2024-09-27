package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Define connection parameters
	connStr := "user=postgres password=tengkorak123 dbname=generic_db sslmode=disable host=localhost"
	
	// Open a connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection by pinging the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	} else {
		fmt.Println("Successfully connected to PostgreSQL!")
	}

	// Execute a sample query
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PostgreSQL Version:", version)
}
