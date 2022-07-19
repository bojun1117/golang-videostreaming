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

func init() {
	connStr := "user=postgres password=eric dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func AddUserCredential(username string, pwd string) error { //新增使用者
	stmtIns, err := db.Prepare("INSERT INTO users (user_name, pwd) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(username, pwd)
	if err != nil {
		return err
	}
	stmtIns.Close()
	return nil
}

func GetUserCredential(username string, pwd string) error { //認證使用者
	stmtOut, err := db.Prepare("SELECT user_id FROM users WHERE user_name=$1 AND pwd=$2")
	if err != nil {
		return err
	}
	var user_id int
	err = stmtOut.QueryRow(username, pwd).Scan(&user_id)
	if err != nil {
		return err
	}
	defer stmtOut.Close()
	return nil
}

func GetUser(user_name string) (*defs.User, error) { //取得使用者資料
	stmtOut, err := db.Prepare("SELECT user_id,pwd FROM users WHERE user_name=$1")
	if err != nil {
		return nil, err
	}
	var user_id int
	var pwd string
	err = stmtOut.QueryRow(user_name).Scan(&user_id, &pwd)
	if err != nil {
		return nil, err
	}
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
	stmtIns, err := db.Prepare("INSERT INTO video_info (author_name, video_title, create_time, viewed) VALUES($1, $2, $3, 0)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(author_name, title, ctime)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func ListVideoInfo() ([]*defs.VideoInfo, error) {
	stmtOut, err := db.Prepare("SELECT * FROM video_info ORDER BY create_time DESC")
	var res []*defs.VideoInfo
	if err != nil {
		return res, err
	}
	rows, err := stmtOut.Query()
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var author, title, ctime string
		var id, view int
		if err := rows.Scan(&id, &author, &title, &ctime, &view); err != nil {
			return res, err
		}
		ctime = ctime[:10]
		vi := &defs.VideoInfo{
			Video_id:    id,
			Author_name: author,
			Video_title: title,
			Create_time: ctime,
			Viewed:      view,
		}
		res = append(res, vi)
	}
	defer stmtOut.Close()
	return res, nil
}

func DeleteVideoInfo(vid int) error {
	stmtDel, err := db.Prepare("DELETE FROM video_info WHERE video_id = $1")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}
