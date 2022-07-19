package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"video-streaming/api/dbops"
	"video-streaming/api/defs"
	"video-streaming/api/session"

	"github.com/julienschmidt/httprouter"
)

func homehandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //首頁
	t, e := template.ParseFiles("./templates/home.html")
	if e != nil {
		log.Printf("Parsing template home.html error: %s", e)
		return
	}

	t.Execute(w, "home.html")
}

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
	err := dbops.GetUserCredential(ubody.Username, ubody.Pwd)
	if err != nil || uname != ubody.Username {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	if session.ValidateUser(w, r, ubody.Username) == true {
		return
	}
	session.RegisterSessionInfo(w, r, ubody.Username)
	return
}

func getUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //取得使用者資訊
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	uname := p.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(u)
	return
}

func addNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //新增影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	err := dbops.AddNewVideo("bojun", "i have a pen")
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	return
}

func listAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //顯示影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	vs, err := dbops.ListVideoInfo()
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	fmt.Println(vs)
}

func deleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) { //刪除影片
	if !session.ValidateUser(w, r, p.ByName("username")) {
		log.Printf("Unauthorized user \n")
		return
	}
	vid, _ := strconv.Atoi(p.ByName("vid"))
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
}
