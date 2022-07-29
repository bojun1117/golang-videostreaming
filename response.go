package main

import (
	"net/http"
	"net/url"
)

func cookieMessage(m string, w http.ResponseWriter) { //回應使用者訊息
	message := m
	http.SetCookie(w, &http.Cookie{
		Name:   "messagecookie",
		Value:  url.QueryEscape(message),
		MaxAge: 1,
	})
}
