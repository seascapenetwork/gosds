// The database package handles all the database operations.
// Note that for now it uses Mysql as a hardcoded data
package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/blocklords/gosds/env"
	_ "github.com/go-sql-driver/mysql"
)

// Creates a DSN (Data Source Name) with the database credentials
func create_dsn() (string, error) {
	if !env.Exists("DB_USER") {
		return "", errors.New("the 'DB_USER' environment variable not set")
	}
	if !env.Exists("DB_PASSWORD") {
		return "", errors.New("the 'DB_PASSWORD' environment variable not set")
	}
	if !env.Exists("DB_HOST") {
		return "", errors.New("the 'DB_HOST' environment variable not set")
	}
	if !env.Exists("DB_PORT") {
		return "", errors.New("the 'DB_PORT' environment variable not set")
	}
	if !env.Exists("DB_NAME") {
		return "", errors.New("the 'DB_NAME' environment variable not set")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		env.GetString("DB_USER"),
		env.GetString("DB_PASSWORD"),
		env.GetString("DB_HOST"),
		env.GetString("DB_PORT"),
		env.GetString("DB_NAME"),
	)
	return dsn, nil
}

// Opens a connection to the database and returns it.
func Open() (*sql.DB, error) {
	dsn, err := create_dsn()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
