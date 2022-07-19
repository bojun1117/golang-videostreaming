package session

import (
	"net/http"
	"video-streaming/api/defs"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
}

func ValidateUser(w http.ResponseWriter, r *http.Request,username string) bool{//確認session存在
	session, err := store.Get(r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	if session.Values["auth"] == true{
		return true
	}
	return false
}

func RegisterSessionInfo(w http.ResponseWriter, r *http.Request,username string) { //註冊session
	session, err := store.Get(r, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_name"] = username
	session.Values["auth"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetSessionInfo(w http.ResponseWriter, r *http.Request,uname string) *defs.SessionInfo { //取得session資訊
	session, err := store.Get(r,uname)
	si := &defs.SessionInfo{
		User_name: session.Values["user_name"].(string),
		Auth: session.Values["auth"].(bool),
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return si
}
