package dbops

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func init() {
	connStr := "user=postgres password=eric dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func AddUserCredential(username string, pwd string) error {
	stmtIns, err := db.Prepare("INSERT INTO users (user_name, pwd) VALUES ($1, $2)")
	CheckErr(err)
	_, err = stmtIns.Exec(username, pwd)
	CheckErr(err)
	stmtIns.Close()
	return nil
}
