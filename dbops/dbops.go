package dbops

import (
	"database/sql"
	"fmt"
	"log"
	"video-streaming/defs"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

const (
	// Initialize connection constants.
	HOST     = "database-18.cqp6xln0cu6w.ap-northeast-1.rds.amazonaws.com"
	DATABASE = "vsproject"
	USER     = "postgres"
	PASSWORD = "eric1117"
)

func init() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)
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

func AddNewVideo(author_name string, title string, cover string) error { //新增影片
	stmtIns, err := db.Prepare("INSERT INTO videos (author_name, video_title, create_time, viewed, cover) VALUES($1, $2, NOW(), 0, $3)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(author_name, title, cover)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func CheckNewVideo(author_name string, title string) bool { //檢查影片是否重複
	stmtIns, err := db.Prepare("SELECT EXISTS(SELECT * FROM videos WHERE author_name=$1 and video_title=$2)")
	if err != nil {
		return true
	}
	var exist bool
	err = stmtIns.QueryRow(author_name, title).Scan(&exist)
	if err != nil {
		return true
	}
	defer stmtIns.Close()
	return exist
}

func GetVideoInfo(vid int) (*defs.VideoInfo, error) { //取得單一影片資訊
	stmtOut, err := db.Prepare("SELECT video_title, viewed FROM videos WHERE video_id=$1")
	if err != nil {
		return nil, err
	}
	var title string
	var view int
	err = stmtOut.QueryRow(vid).Scan(&title, &view)
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{
		Video_id:    vid,
		Video_title: title,
		Viewed:      view,
	}
	return res, nil
}

func ListVideoInfo(username string) ([]*defs.VideoInfo, error) { //顯示影片
	var res []*defs.VideoInfo
	var rows *sql.Rows
	var stmtOut *sql.Stmt
	if username == "" { //所有
		stmtOut, err = db.Prepare("SELECT * FROM videos ORDER BY viewed DESC")
		if err != nil {
			return res, err
		}
		rows, err = stmtOut.Query()
		if err != nil {
			return res, err
		}
	} else { //僅使用者上傳
		stmtOut, err = db.Prepare("SELECT * FROM videos WHERE author_name=$1 ORDER BY viewed DESC")
		if err != nil {
			return res, err
		}
		rows, err = stmtOut.Query(username)
		if err != nil {
			return res, err
		}
	}
	for rows.Next() {
		var author, title, ctime, cover string
		var id, view int
		if err := rows.Scan(&id, &author, &title, &ctime, &view, &cover); err != nil {
			return res, err
		}
		ctime = ctime[:10]
		vi := &defs.VideoInfo{
			Video_id:    id,
			Author_name: author,
			Video_title: title,
			Create_time: ctime,
			Viewed:      view,
			Cover:       cover,
		}
		res = append(res, vi)
	}
	defer stmtOut.Close()
	return res, nil
}

func ListSpecifyVideos(q string) ([]*defs.VideoInfo, error) { //搜尋特定結果
	var res []*defs.VideoInfo
	stmtOut, err := db.Prepare("SELECT * FROM videos WHERE video_title LIKE $1 ORDER BY viewed DESC limit 24")
	if err != nil {
		return res, err
	}
	rows, err := stmtOut.Query("%" + q + "%")
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var author, title, ctime, cover string
		var id, view int
		if err := rows.Scan(&id, &author, &title, &ctime, &view, &cover); err != nil {
			return res, err
		}
		ctime = ctime[:10]
		vi := &defs.VideoInfo{
			Video_id:    id,
			Author_name: author,
			Video_title: title,
			Create_time: ctime,
			Viewed:      view,
			Cover:       cover,
		}
		res = append(res, vi)
	}
	defer stmtOut.Close()
	return res, nil
}

func DeleteVideoInfo(vid int, uname string) error { //刪除影片
	stmtDelv, err := db.Prepare("DELETE FROM videos WHERE video_id = $1 and author_name = $2")
	if err != nil {
		return err
	}
	_, err = stmtDelv.Exec(vid, uname)
	if err != nil {
		return err
	}
	defer stmtDelv.Close()
	stmtDelc, err := db.Prepare("DELETE FROM comments WHERE video_id = $1")
	if err != nil {
		return err
	}
	_, err = stmtDelc.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDelc.Close()
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
	var res []*defs.Comment
	stmtOut, err := db.Prepare("SELECT * from comments where video_id=$1")
	if err != nil {
		return res, err
	}
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
		tmp := []byte(time)
		tmp[10] = ' '
		time = string(tmp)
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

func AddViewCount(vid int, viewed int) error { //點率加一
	stmout, err := db.Prepare("Update videos SET viewed=$1 WHERE video_id=$2")
	if err != nil {
		return err
	}
	viewed++
	_, err = stmout.Exec(viewed, vid)
	if err != nil {
		return err
	}
	defer stmout.Close()
	return nil
}
