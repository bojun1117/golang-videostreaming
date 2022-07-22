package defs

//requests
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewComment struct {
	User_name string `json:"user_name"`
	Content   string `json:"contents"`
}

type NewVideo struct {
	Author string `json:"author"`
	Title  string `json:"title"`
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
	Video_id    int
	Author_name string `json:"author_name"`
	Video_title string `json:"video_title"`
	Create_time string `json:"create_time"`
	Viewed      int    `json:"Viewed"`
}

type Comment struct {
	Comment_id  int
	User_name   string `json:"user_name"`
	Content     string `json:"content"`
	Video_id    int    `json:"video_id"`
	Record_time string `json:"record_time"`
}

type SessionInfo struct {
	User_name string
	Auth      bool
}
