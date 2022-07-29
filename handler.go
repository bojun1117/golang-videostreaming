package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"video-streaming/dbops"
	"video-streaming/defs"
	"video-streaming/session"

	"github.com/julienschmidt/httprouter"
)

//templates
const TEMPLATE_DIR = "./webserver/templates/"

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //首頁
	var message string
	em, err := r.Cookie("messagecookie")
	if err != nil {
		message = ""
	} else {
		message, _ = url.QueryUnescape(em.Value)
	}

	user := session.ValidateUser(w, r)

	t, e := template.ParseFiles(TEMPLATE_DIR + "home.html")
	if e != nil {
		log.Printf("Parsing template home.html error: %s", e)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" { //無搜尋
		vs, err := dbops.ListVideoInfo("")
		if err != nil {
			return
		}
		vsi := &defs.VideosInfo{
			Videos:  vs,
			User:    user,
			Message: message,
		}
		t.Execute(w, vsi)
		return
	}
	vs, err := dbops.ListSpecifyVideos(query) //搜尋
	if err != nil {
		return
	}
	vsi := &defs.VideosInfo{
		Videos:  vs,
		User:    user,
		Message: message,
	}

	t.Execute(w, vsi)
}

func createUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //創建使用者頁面
	var message string
	em, err := r.Cookie("messagecookie")
	if err != nil {
		message = ""
	} else {
		message, _ = url.QueryUnescape(em.Value)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "adduser.html")
	if e != nil {
		log.Printf("Parsing template adduser.html error: %s", e)
		return
	}

	t.Execute(w, message)
}

func login(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //登入頁面
	var message string
	em, err := r.Cookie("messagecookie")
	if err != nil {
		message = ""
	} else {
		message, _ = url.QueryUnescape(em.Value)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "login.html")
	if e != nil {
		log.Printf("Parsing template login.html error: %s", e)
		return
	}

	t.Execute(w, message)
}

func userVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params){ //會員頁面
	user := session.ValidateUser(w, r)
	if user == ""{
		http.Redirect(w, r, "./videos", http.StatusFound)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "uservideos.html")
	if e != nil {
		log.Printf("Parsing template uservideos.html error: %s", e)
		return
	}

	vs,err := dbops.ListVideoInfo(user)
	if err != nil {
		return
	}
	vsi := &defs.VideosInfo{
		Videos: vs,
		User: user,
	}

	t.Execute(w, vsi)
}

//database
func userInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //創建使用者
	ubody := &defs.UserCredential{
		Username: r.PostFormValue("user"),
		Pwd:      r.PostFormValue("password"),
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		var message string
		if err.Error() == "blank"{
			message = "不可空白"
		}else{
			message = "使用者已被註冊"
		}
		cookieMessage(message, w)
		http.Redirect(w, r, "./user", http.StatusFound)
		return
	}
	message := "註冊成功"
	cookieMessage(message, w)

	http.Redirect(w, r, "./videos", http.StatusFound)
}

func loginCredential(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //登入驗證
	Username := r.PostFormValue("user")
	Pwd := r.PostFormValue("password")

	user_id, err := dbops.GetUserCredential(Username, Pwd)
	if err != nil {
		message := "帳號或密碼錯誤"
		cookieMessage(message, w)
		http.Redirect(w, r, "./login", http.StatusFound)
		return
	}
	session.RegisterSessionInfo(w, r, Username, user_id)
	http.Redirect(w, r, "./videos", http.StatusFound)
	return
}

func logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //登出
	session.RemoveSessionAuth(w, r)
	http.Redirect(w, r, "./videos", http.StatusFound)
}

func getUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //取得使用者資訊

	uname := p.ByName("username")
	ubody, err := dbops.GetUser(uname)
	if err != nil {

		return
	}
	fmt.Println(ubody)
	return
}

func addNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //新增影片

	uname := p.ByName("username")
	vbody := defs.NewVideo{
		Author: uname,
		Title:  "i have a pen",
	}
	err := dbops.AddNewVideo(vbody.Author, vbody.Title)
	if err != nil {

		return
	}
	return
}

func listAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示所有影片

	vs, err := dbops.ListVideoInfo("")
	if err != nil {

		return
	}
	fmt.Println(vs[0].Create_time)
}

func listUserVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示使用者影片

	username := p.ByName("username")
	vs, err := dbops.ListVideoInfo(username)
	if err != nil {

		return
	}
	fmt.Println(vs)
}

func getVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //取得影片資訊

	vid, _ := strconv.Atoi(p.ByName("vid"))
	vbody, err := dbops.GetVideoInfo(vid)
	if err != nil {

		return
	}
	fmt.Println(vbody.Video_title)
	return
}

func deleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除影片

	vid, _ := strconv.Atoi(p.ByName("vid"))
	uname := p.ByName("username")
	err := dbops.DeleteVideoInfo(vid, uname)
	if err != nil {

		return
	}
}

func postComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //新增評論

	uname := p.ByName("username")
	cbody := &defs.NewComment{
		User_name: uname,
		Content:   "test",
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	if err := dbops.AddNewComments(vid, cbody.User_name, cbody.Content); err != nil {
		log.Printf("Error in PostComment: %s", err)

	}
	return
}

func showComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示評論

	vid, _ := strconv.Atoi(p.ByName("vid"))
	cm, err := dbops.ListComments(vid)
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)

		return
	}
	fmt.Println(cm[0].Record_time)
}

func deleteComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除影片

	cid, _ := strconv.Atoi(p.ByName("cid"))
	uname := p.ByName("username")
	err := dbops.DeleteCommentInfo(cid, uname)
	if err != nil {

		return
	}
}

//streaming
const (
	VIDEO_DIR       = "./videos/"
	MAX_UPLOAD_SIZE = 50 * 1024 * 1024 // 50MB
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid")
	vl := VIDEO_DIR + vid
	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Error when try to open file: %v", err)

		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {

		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {

		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)

	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)

		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload successfully")
}
