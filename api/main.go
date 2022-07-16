package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler { //初始化handler
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //實作ServeHTTP method
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router { //API控制
	router := httprouter.New()
	router.GET("/", homehandler)
	router.GET("/user", createUser)
	router.GET("/user/:username", login)
	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))
	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":5050", mh)
}
