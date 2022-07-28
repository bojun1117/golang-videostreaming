package main

import (
	"net/http"
	"net/url"
)

func cookieMessage(m string, w http.ResponseWriter) {
	message := m
	http.SetCookie(w, &http.Cookie{
		Name:   "messagecookie",
		Value:  url.QueryEscape(message),
		MaxAge: 1,
	})

}
