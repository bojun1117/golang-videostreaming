package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func homehandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {		//首頁
	t, e := template.ParseFiles("./templates/home.html")
	if e != nil {
		log.Printf("Parsing template home.html error: %s", e)
		return
	}

	t.Execute(w, "home.html")
}

func login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
