package db

import (
	"aspire/util"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	HOST = "localhost"
	PORT = "5432"
)

type DB struct {
	*sql.DB
}

func Init() (*DB, error) {
	username := util.GetEnvOrDefaultString("DB_USERNAME", "avinragh")
	password := util.GetEnvOrDefaultString("DB_PASSWORD", "toor")
	database := util.GetEnvOrDefaultString("DB_DATABASE", "aspire")

	connectionString := fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable",
		PORT, HOST, username, password, database)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Database connection established")
	return &DB{db}, nil
}
