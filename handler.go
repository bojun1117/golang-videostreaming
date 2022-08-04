package main

import (
	"html/template"
	"io"
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

	t, e := template.ParseFiles(TEMPLATE_DIR + "user.html")
	if e != nil {
		log.Printf("Parsing template user.html error: %s", e)
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
	vid, _ := strconv.Atoi(p.ByName("vid"))
	vbody, err := dbops.GetVideoInfo(vid)
	if err != nil {
		log.Printf("Error in getvideoinfo: %s", err)
	}

	streamHandler(w, r, vbody.Video_title)
}

func commentInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //評論頁面
	user := session.ValidateUser(w, r)
	var message string
	em, err := r.Cookie("messagecookie")
	if err == nil {
		message, _ = url.QueryUnescape(em.Value)
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	cbody, err := dbops.ListComments(vid)
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
	}
	ci := &defs.CommentsInfo{
		Comments: cbody,
		User:     user,
		Message:  message,
	}

	t, e := template.ParseFiles(TEMPLATE_DIR + "comment.html")
	if e != nil {
		log.Printf("Parsing template comment.html error: %s", e)
		return
	}

	t.Execute(w, ci)
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
		message := "請先登入"
		cookieMessage(message, w)
		http.Redirect(w, r, h, http.StatusFound)
		return
	}
	content := r.PostFormValue("content")
	if content == "" {
		message := "不可空白"
		cookieMessage(message, w)
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
	cover := r.PostFormValue("cover")
	user := session.ValidateUser(w, r)
	err := dbops.AddNewVideo(user, title, cover)
	if err != nil {
		message := "影片名稱重複"
		cookieMessage(message, w)
		http.Redirect(w, r, "upload", http.StatusFound)
		return
	}

	r.ParseMultipartForm(50)
	video, _, err := r.FormFile("video")
	if err != nil {
		message := "影片上傳失敗"
		cookieMessage(message, w)
		http.Redirect(w, r, "user", http.StatusFound)
		return
	}
	defer video.Close()
	v, err := os.OpenFile("videos/"+title+".mp4", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		message := "影片上傳失敗"
		cookieMessage(message, w)
		http.Redirect(w, r, "user", http.StatusFound)
		return
	}
	defer v.Close()
	io.Copy(v, video)

	http.Redirect(w, r, "user", http.StatusFound)
}

//streaming
const VIDEO_DIR = "./videos/"

func streamHandler(w http.ResponseWriter, r *http.Request, vname string) {
	vpath := VIDEO_DIR + vname + ".mp4"
	video, err := os.Open(vpath)
	if err != nil {
		log.Printf("Error when try to open file: %v", err)
		return
	}

	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}
