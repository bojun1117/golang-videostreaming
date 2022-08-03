package main

import (
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
const TEMPLATE_DIR = "./templates/"

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //首頁
	var message string
	em, err := r.Cookie("messagecookie")
	if err == nil {
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

func createUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //註冊頁面
	var message string
	em, err := r.Cookie("messagecookie")
	if err == nil {
		message, _ = url.QueryUnescape(em.Value)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "register.html")
	if e != nil {
		log.Printf("Parsing template register.html error: %s", e)
		return
	}

	t.Execute(w, message)
}

func login(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //登入頁面
	var message string
	em, err := r.Cookie("messagecookie")
	if err == nil {
		message, _ = url.QueryUnescape(em.Value)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "login.html")
	if e != nil {
		log.Printf("Parsing template login.html error: %s", e)
		return
	}

	t.Execute(w, message)
}

func userVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //會員頁面
	user := session.ValidateUser(w, r)
	if user == "" {
		log.Printf("not member")
		http.Redirect(w, r, "videos", http.StatusFound)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "uservideos.html")
	if e != nil {
		log.Printf("Parsing template uservideos.html error: %s", e)
		return
	}

	vs, err := dbops.ListVideoInfo(user)
	if err != nil {
		return
	}
	vsi := &defs.VideosInfo{
		Videos: vs,
		User:   user,
	}

	t.Execute(w, vsi)
}

func videoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //單一影片頁面
	user := session.ValidateUser(w, r)
	vid, _ := strconv.Atoi(p.ByName("vid"))
	vbody, err := dbops.GetVideoInfo(vid)
	if err != nil {
		log.Printf("Error in getvideoinfo: %s", err)
	}
	cbody, err := dbops.ListComments(vid)
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
	}
	vdi := &defs.VideoDetails{
		Comments: cbody,
		Title:    vbody.Video_title,
		Author:   vbody.Author_name,
		User:     user,
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "video.html")
	if e != nil {
		log.Printf("Parsing template video.html error: %s", e)
		return
	}

	t.Execute(w, vdi)
}

func upload(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //上傳頁面
	user := session.ValidateUser(w, r)
	if user == "" {
		log.Printf("not member")
		http.Redirect(w, r, "videos", http.StatusFound)
	}

	var message string
	em, err := r.Cookie("messagecookie")
	if err == nil {
		message, _ = url.QueryUnescape(em.Value)
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "upload.html")
	if e != nil {
		log.Printf("Parsing template upload.html error: %s", e)
		return
	}

	t.Execute(w, message)
}

//database
func userInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //註冊
	username := r.PostFormValue("user")
	pwd := r.PostFormValue("password")
	if err := dbops.AddUserCredential(username, pwd); err != nil {
		message := "使用者已被註冊"
		cookieMessage(message, w)
		http.Redirect(w, r, "./create", http.StatusFound)
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

func deleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除影片
	user := session.ValidateUser(w, r)
	if user == "" {
		return
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	err := dbops.DeleteVideoInfo(vid, user)
	if err != nil {
		log.Printf("error")
	}
	http.Redirect(w, r, "../user", http.StatusFound)
}

func postComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //新增評論
	user := session.ValidateUser(w, r)
	h := p.ByName("vid")
	if user == "" {
		log.Printf("not a user")
		http.Redirect(w, r, h, http.StatusFound)
		return
	}
	content := r.PostFormValue("content")
	if content == "" {
		log.Printf("blank")
		http.Redirect(w, r, h, http.StatusFound)
		return
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	if err := dbops.AddNewComments(vid, user, content); err != nil {
		log.Printf("Error in PostComment: %s", err)
	}
	http.Redirect(w, r, h, http.StatusFound)
}

func deleteComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除評論
	h := "../" + p.ByName("vid")
	cid, _ := strconv.Atoi(p.ByName("cid"))
	user := session.ValidateUser(w, r)
	err := dbops.DeleteCommentInfo(cid, user)
	if err != nil {
		log.Printf("Error in deleteComment: %s", err)
	}
	http.Redirect(w, r, h, http.StatusFound)
}

func uploadVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //上傳影片
	title := r.PostFormValue("title")
	user := session.ValidateUser(w, r)
	err := dbops.AddNewVideo(user, title)
	if err != nil {
		message := "影片名稱重複"
		cookieMessage(message, w)
		http.Redirect(w, r, "upload", http.StatusFound)
		return
	}

	r.ParseMultipartForm(50)
	video, handler, err := r.FormFile("video")
	if err != nil {
		message := "影片上傳失敗"
		cookieMessage(message, w)
		http.Redirect(w, r, "user", http.StatusFound)
		return
	}
	defer video.Close()
	v, err := os.OpenFile("videos/"+handler.Filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		message := "影片上傳失敗"
		cookieMessage(message, w)
		http.Redirect(w, r, "user", http.StatusFound)
		return
	}
	defer v.Close()
	io.Copy(v, video)

	cover, handler, err := r.FormFile("cover")
	if err != nil {
		message := "封面上傳失敗"
		cookieMessage(message, w)
		http.Redirect(w, r, "user", http.StatusFound)
		return
	}
	defer video.Close()
	c, err := os.OpenFile("videos/"+handler.Filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		message := "封面上傳失敗"
		cookieMessage(message, w)
		http.Redirect(w, r, "user", http.StatusFound)
		return
	}
	defer c.Close()
	io.Copy(c, cover)

	http.Redirect(w, r, "user", http.StatusFound)
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
