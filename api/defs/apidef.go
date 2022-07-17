package defs

//requests
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"contents"`
}

type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

// response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type UserSession struct {
	UserID  int `json:"user_id"`
	SessionId string `json:"session_id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type Comments struct {
	Comments []*Comment `json:"comments"`
}

// Data model

type User struct {
	User_id   int
	User_name string
	Pwd       string
}

type VideoInfo struct {
	Video_id      int
	Author_name   string `json:"author_name"`
	Video_title   string `json:"video_title"`
	Display_ctime string `json:"display_ctime"`
	Create_time   string `json:"create_time"`
}

type Comment struct {
	Comment_id  int 
	Video_title string `json:"video_title"`
	User_name   string `json:"user_name"`
	Content     string `json:"content"`
	Record_time string `json:"record_time"`
}

type SimpleSession struct {
	Session_id int
	TTL        int64
	User_id    int
}
