package session

import (
	"net/http"
	"video-streaming/main/dbops"
	"video-streaming/main/defs"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
}

func ValidateUser(w http.ResponseWriter, r *http.Request, username string) bool { //確認session存在
	session, err := store.Get(r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	user, err := dbops.GetUser(username)
	if err != nil {
		return false
	}
	if session.Values["auth"] == true && session.Values["user_id"] == user.User_id {
		return true
	}
	return false
}

func RegisterSessionInfo(w http.ResponseWriter, r *http.Request, username string, user_id int) { //註冊session
	session, err := store.Get(r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = user_id
	session.Values["auth"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetSessionInfo(w http.ResponseWriter, r *http.Request, uname string) *defs.SessionInfo { //取得session資訊
	session, err := store.Get(r, uname)
	si := &defs.SessionInfo{
		User_id: session.Values["user_id"].(int),
		Auth:    session.Values["auth"].(bool),
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return si
}
