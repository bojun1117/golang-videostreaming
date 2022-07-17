package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"video-streaming/api/dbops"
	"video-streaming/api/defs"

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
	user_id,err := dbops.GetUserCredential(ubody.Username, ubody.Pwd)
	if err != nil || uname != ubody.Username{
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}
	ui := &defs.UserSession{
		UserID: user_id,
	}
	_ = ui
	return
}

func getUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params){		//取得使用者資訊
	u, err := dbops.GetUser(3)
	if err != nil {
		log.Printf("Erorr in GetUserinfo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	_ = u
	return
}