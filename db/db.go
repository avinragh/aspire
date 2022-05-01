package db

import (
	"aspire/constants"
	"aspire/util"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Init() (*DB, error) {
	username := util.GetEnvOrDefaultString(constants.EnvDbUsername, constants.DbUsername)
	password := util.GetEnvOrDefaultString(constants.EnvDbPassword, constants.DbPassword)
	database := util.GetEnvOrDefaultString(constants.EnvDbDatabase, constants.DbDatabase)
	host := util.GetEnvOrDefaultString(constants.EnvDbHost, constants.DbHost)
	port := util.GetEnvOrDefaultString(constants.EnvDbPort, constants.DbPort)

	connectionString := fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=disable",
		port, host, username, password, database)
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
