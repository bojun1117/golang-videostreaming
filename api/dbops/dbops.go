package dbops

import (
	"database/sql"
	"video-streaming/api/defs"

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

func AddUserCredential(username string, pwd string) error { //新增使用者
	stmtIns, err := db.Prepare("INSERT INTO users (user_name, pwd) VALUES ($1, $2)")
	CheckErr(err)
	_, err = stmtIns.Exec(username, pwd)
	CheckErr(err)
	stmtIns.Close()
	return nil
}

func GetUserCredential(username string, pwd string) (int, error) { //認證使用者
	stmtOut, err := db.Prepare("SELECT user_id FROM users WHERE user_name=$1 AND pwd=$2")
	CheckErr(err)
	var user_id int
	err = stmtOut.QueryRow(username, pwd).Scan(&user_id)
	CheckErr(err)
	defer stmtOut.Close()
	return user_id, nil
}

func GetUser(user_id int) (*defs.User, error) {		//取得使用者資料
	stmtOut, err := db.Prepare("SELECT user_name,pwd FROM users WHERE user_id=$1")
	CheckErr(err)
	var user_name string
	var pwd string
	err = stmtOut.QueryRow(user_id).Scan(&user_name, &pwd)
	CheckErr(err)
	res := &defs.User{
		User_id:   user_id,
		User_name: user_name,
		Pwd:       pwd,
	}
	defer stmtOut.Close()
	return res, nil
}
