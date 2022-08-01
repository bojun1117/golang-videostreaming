package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler { //初始化handler
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //實作ServeHTTP method
	if !m.l.GetConn() {
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router { //API控制
	router := httprouter.New()
	router.GET("/videos", homeHandler)
	router.GET("/create", createUser)
	router.POST("/create", userInfo)
	router.GET("/login", login)
	router.POST("/login", loginCredential)
	router.GET("/logout", logout)
	router.GET("/user", userVideos)
	router.GET("/user/:vid", deleteVideo)
	router.GET("/videos/:vid", videoInfo)
	/*

		router.POST("/videos/:vid",postComment)
		router.GET("/videos/:vid/:cid",deleteComment)
		router.GET("/videos/:username/upload",upload)
		router.POST("/videos/:username/upload",uploadVideo)
	*/
	router.ServeFiles("/css/*filepath", http.Dir("./webserver/css"))
	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 10)
	http.ListenAndServe(":5050", mh)
}
