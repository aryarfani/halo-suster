package db

import (
	"database/sql"
	"eniqilo-store/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Database instance
var DB *sql.DB

// Database settings
var (
	host     = config.GetConfig("DB_HOST")
	port     = config.GetConfig("DB_PORT")
	user     = config.GetConfig("DB_USERNAME")
	password = config.GetConfig("DB_PASSWORD")
	dbname   = config.GetConfig("DB_NAME")
	dbparams = config.GetConfig("DB_PARAMS")
)

// Connect function
func Connect() error {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		user,
		password,
		host,
		port,
		dbname,
		dbparams,
	)

	fmt.Println(connStr)
	DB, err = sql.Open("postgres", connStr)

	DB.SetMaxOpenConns(20)
	DB.SetMaxIdleConns(10)

	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}

	return nil
}

// perform to close mysql database connection
func Close() error {
	if err := DB.Close(); err != nil {
		log.Fatalf("error when closing the database connection: %v", err)
		return err
	}

	log.Println("database connection closed")

	return nil
}
