package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router {
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
	router.GET("/comments/:vid", commentInfo)
	router.POST("/comments/:vid", postComment)
	router.GET("/comments/:vid/:cid", deleteComment)
	router.GET("/upload", upload)
	router.POST("/upload", uploadVideo)
	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":5000", mh)
}
