package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

func homehandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //首頁
	t, e := template.ParseFiles(TEMPLATE_DIR + "index.html")
	if e != nil {
		log.Printf("Parsing template home.html error: %s", e)
		return
	}
	t.Execute(w, "home.html")
}

//database
func createUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //創建使用者
	res := []byte(`{
		"user_name":"bojun",
		"pwd":"eric1117"
	}`)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //登入
	ubody := &defs.UserCredential{
		Username: "bojun",
		Pwd:      "eric1117",
	}
	uname := p.ByName("username")
	user_id, err := dbops.GetUserCredential(ubody.Username, ubody.Pwd)
	if err != nil || uname != ubody.Username {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	if session.ValidateUser(w, r, ubody.Username) == true {
		return
	}
	session.RegisterSessionInfo(w, r, ubody.Username, user_id)
	return
}

func getUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //取得使用者資訊
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	uname := p.ByName("username")
	ubody, err := dbops.GetUser(uname)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(ubody)
	return
}

func addNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //新增影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	uname := p.ByName("username")
	vbody := defs.NewVideo{
		Author: uname,
		Title:  "i have a pen",
	}
	err := dbops.AddNewVideo(vbody.Author, vbody.Title)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	return
}

func listAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示所有影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	vs, err := dbops.ListVideoInfo("")
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(vs[0].Create_time)
}

func listUserVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示使用者影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	username := p.ByName("username")
	vs, err := dbops.ListVideoInfo(username)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(vs)
}

func getVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //取得影片資訊
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	vbody, err := dbops.GetVideoInfo(vid)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(vbody.Video_title)
	return
}

func deleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	uname := p.ByName("username")
	err := dbops.DeleteVideoInfo(vid, uname)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
}

func postComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //新增評論
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	uname := p.ByName("username")
	cbody := &defs.NewComment{
		User_name: uname,
		Content:   "test",
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	if err := dbops.AddNewComments(vid, cbody.User_name, cbody.Content); err != nil {
		log.Printf("Error in PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
	}
	return
}

func showComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示評論
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	cm, err := dbops.ListComments(vid)
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(cm[0].Record_time)
}

func deleteComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	cid, _ := strconv.Atoi(p.ByName("cid"))
	uname := p.ByName("username")
	err := dbops.DeleteCommentInfo(cid, uname)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
}

//streaming
const (
	VIDEO_DIR = "./videos/"
	MAX_UPLOAD_SIZE = 50 * 1024 * 1024 // 50MB
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid")
	vl := VIDEO_DIR + vid
	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Error when try to open file: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload successfully")
}
