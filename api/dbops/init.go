package dbops

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func init() {
	connStr := "user=postgres password=eric dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}
