package main

import (
	"net/http"
	"video-streaming/main/defs"

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
		sendErrorResponse(w, defs.ErrorTooManyRequests)
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func RegisterHandlers() *httprouter.Router { //API控制
	router := httprouter.New()
	router.GET("/", homehandler)
	router.GET("/user", createUser)
	router.POST("/user/:username", login)
	router.GET("/user/:username", getUserInfo)
	router.POST("/user/:username/videos", addNewVideo)
	router.GET("/user/:username/videos", listAllVideos)
	router.GET("/user/:username/videos/:vid", getVideo)
	router.DELETE("/user/:username/videos/:vid", deleteVideo)
	router.POST("/user/:username/videos/:vid/comments", postComment)
	router.GET("/user/:username/videos/:vid/comments", showComments)
	router.DELETE("/user/:username/videos/:vid/comments/:cid", deleteComment)
	router.GET("/videos/:vid", streamHandler)
	router.POST("/upload/:vid", uploadHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))
	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 10)
	http.ListenAndServe(":5050", mh)
}
