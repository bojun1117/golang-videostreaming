package dbops

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

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

func GetUserCredential(username string, pwd string) (int, error) {
	stmtOut, err := db.Prepare("SELECT user_id FROM users WHERE user_name=$1 AND pwd=$2")
	CheckErr(err)
	var userID int
	err = stmtOut.QueryRow(username, pwd).Scan(&userID)
	CheckErr(err)
	defer stmtOut.Close()
	return userID, nil
}
