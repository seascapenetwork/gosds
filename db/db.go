package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func EnvVars() []string {
	return []string{
		"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME",
	}
}

func getConf(name string) string {
	value := os.Getenv("DB_" + name)
	if len(value) == 0 {
		panic("no 'DB_" + name + "' environment variable set")
	}
	return value
}

func Dsn() string {
	dsn := getConf("USER") + ":" + getConf("PASSWORD") + "@tcp(" + getConf("HOST") + ":" + getConf("PORT") + ")/" + getConf("NAME")
	return dsn
}

func Open() *sql.DB {
	db, err := sql.Open("mysql", Dsn())
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	return db
}
