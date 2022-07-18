package dbops

import (
	"database/sql"
	"log"
	"time"
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

func init() {
	connStr := "user=postgres password=eric dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
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

func GetUserCredential(username string, pwd string) error { //認證使用者
	stmtOut, err := db.Prepare("SELECT user_id FROM users WHERE user_name=$1 AND pwd=$2")
	CheckErr(err)
	var user_id int
	err = stmtOut.QueryRow(username, pwd).Scan(&user_id)
	CheckErr(err)
	defer stmtOut.Close()
	return nil
}

func GetUser(user_name string) (*defs.User, error) { //取得使用者資料
	stmtOut, err := db.Prepare("SELECT user_id,pwd FROM users WHERE user_name=$1")
	CheckErr(err)
	var user_id int
	var pwd string
	err = stmtOut.QueryRow(user_name).Scan(&user_id, &pwd)
	CheckErr(err)
	res := &defs.User{
		User_id:   user_id,
		User_name: user_name,
		Pwd:       pwd,
	}
	defer stmtOut.Close()
	return res, nil
}

func AddNewVideo(author_name string, title string) error {
	t := time.Now()
	ctime := t.Format("2006-01-02") // YY-MM-DD
	stmtIns, err := db.Prepare("INSERT INTO video_info (author_name, video_title, create_time) VALUES($1, $2, $3)")
	CheckErr(err)
	_, err = stmtIns.Exec(author_name, title, ctime)
	CheckErr(err)
	defer stmtIns.Close()
	return nil
}
