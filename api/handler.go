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

func createUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
