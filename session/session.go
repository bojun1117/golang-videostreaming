package session

import (
	"net/http"
	"video-streaming/dbops"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var sessionCookieName = "user-session"

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
}

func ValidateUser(w http.ResponseWriter, r *http.Request) string { //確認session存在 存在返回user
	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		return ""
	}
	if session.Values["user_name"] == nil {
		return ""
	}
	user, err := dbops.GetUser(session.Values["user_name"].(string))
	if err != nil {
		return ""
	}
	if session.Values["auth"] == true {
		return user.User_name
	}
	return ""
}

func RegisterSessionInfo(w http.ResponseWriter, r *http.Request, username string, userid int) { //註冊session
	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = userid
	session.Values["user_name"] = username
	session.Values["auth"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RemoveSessionAuth(w http.ResponseWriter, r *http.Request) { //取消session auth
	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		return
	}
	session.Values["auth"] = false
	err = session.Save(r, w)
	if err != nil {
		return
	}
}
