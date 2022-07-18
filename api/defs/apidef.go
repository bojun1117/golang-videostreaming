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
	Create_time   string `json:"create_time"`
}

type Comment struct {
	Comment_id  int
	Video_title string `json:"video_title"`
	User_name   string `json:"user_name"`
	Content     string `json:"content"`
	Record_time string `json:"record_time"`
}

type SessionInfo struct {
	User_name string
	Auth      bool
}
