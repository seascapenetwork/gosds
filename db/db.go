package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

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

func CreateStmt(db *sql.DB) *sql.Stmt {
	// Prepare statement for reading data
	createStmt, err := db.Prepare(`INSERT IGNORE INTO categorizer_transactions 
	(network_id, address, block_number, txid, tx_index, tx_from, method_name, args, value)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer createStmt.Close()

	return createStmt
}

func SetBlockNumberStmt(db *sql.DB) *sql.Stmt {
	// Prepare statement for reading data
	createStmt, err := db.Prepare(`UPDATE categorizer_blocks SET synced_block = ? WHERE network_id = ? AND address = ? `)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer createStmt.Close()

	return createStmt
}
