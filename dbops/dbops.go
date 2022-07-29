package dbops

import (
	"database/sql"
	"errors"
	"log"
	"video-streaming/defs"

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
	if username == "" || pwd == "" {
		err := errors.New("blank")
		return err
	}
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

func GetUserCredential(username string, pwd string) (int, error) { //認證使用者
	stmtOut, err := db.Prepare("SELECT user_id FROM users WHERE user_name=$1 AND pwd=$2")
	if err != nil {
		return 0, err
	}
	var user_id int
	err = stmtOut.QueryRow(username, pwd).Scan(&user_id)
	if err != nil {
		return 0, err
	}
	defer stmtOut.Close()
	return user_id, nil
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

func AddNewVideo(author_name string, title string) error { //新增影片
	stmtIns, err := db.Prepare("INSERT INTO video_info (author_name, video_title, create_time, viewed) VALUES($1, $2, NOW(), 0)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(author_name, title)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func GetVideoInfo(vid int) (*defs.VideoInfo, error) { //取得單一影片資訊
	stmtOut, err := db.Prepare("SELECT author_name, video_title, create_time, viewed FROM video_info WHERE video_id=$1")
	if err != nil {
		return nil, err
	}
	var author, title, ctime string
	var view int
	err = stmtOut.QueryRow(vid).Scan(&author, &title, &ctime, &view)
	res := &defs.VideoInfo{
		Video_id:    vid,
		Author_name: author,
		Video_title: title,
		Create_time: ctime,
		Viewed:      view,
	}
	return res, nil
}

func ListVideoInfo(username string) ([]*defs.VideoInfo, error) { //顯示影片
	var res []*defs.VideoInfo
	var rows *sql.Rows
	var stmtOut *sql.Stmt
	if username == "" { //所有
		stmtOut, err = db.Prepare("SELECT * FROM video_info ORDER BY create_time DESC")
		if err != nil {
			return res, err
		}
		rows, err = stmtOut.Query()
		if err != nil {
			return res, err
		}
	} else { //僅使用者上傳
		stmtOut, err = db.Prepare("SELECT * FROM video_info WHERE author_name=$1 ORDER BY create_time DESC")
		if err != nil {
			return res, err
		}
		rows, err = stmtOut.Query(username)
		if err != nil {
			return res, err
		}
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

func ListSpecifyVideos(q string) ([]*defs.VideoInfo, error) { //搜尋特定結果
	var res []*defs.VideoInfo
	stmtOut, err := db.Prepare("SELECT * FROM video_info WHERE video_title LIKE $1")
	if err != nil {
		return res, err
	}
	rows, err := stmtOut.Query("%" + q + "%")
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

func DeleteVideoInfo(vid int, uname string) error { //刪除影片
	stmtDel, err := db.Prepare("DELETE FROM video_info WHERE video_id = $1 and author_name = $2")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid, uname)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

func AddNewComments(vid int, user_name string, content string) error { //新增評論
	stmtIns, err := db.Prepare("INSERT INTO comments (user_name, video_id, contents, record_time) VALUES ($1, $2, $3, NOW())")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(user_name, vid, content)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func ListComments(vid int) ([]*defs.Comment, error) { //顯示評論
	stmtOut, err := db.Prepare("SELECT * from comments where video_id=$1")
	var res []*defs.Comment
	rows, err := stmtOut.Query(vid)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var name, content, time string
		var id int
		if err := rows.Scan(&id, &name, &content, &vid, &time); err != nil {
			return res, err
		}
		c := &defs.Comment{
			Comment_id:  id,
			User_name:   name,
			Content:     content,
			Video_id:    vid,
			Record_time: time[:19],
		}
		res = append(res, c)
	}
	defer stmtOut.Close()
	return res, nil
}

func DeleteCommentInfo(cid int, uname string) error { //刪除評論
	stmtDel, err := db.Prepare("DELETE FROM comments WHERE comment_id = $1 and user_name = $2")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(cid, uname)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}
