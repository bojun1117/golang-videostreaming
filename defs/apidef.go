package defs

// response
type VideosInfo struct {
	Videos  []*VideoInfo
	User    string
	Message string
}

type VideoDetails struct {
	Comments []*Comment
	Title    string
	Author   string
	User     string
}

// Data model

type User struct {
	User_id   int
	User_name string
	Pwd       string
}

type VideoInfo struct {
	Video_id    int
	Author_name string
	Video_title string
	Create_time string
	Viewed      int
}

type Comment struct {
	Comment_id  int
	User_name   string
	Content     string
	Video_id    int
	Record_time string
}

type SessionInfo struct {
	User_id   int
	User_name string
	Auth      bool
}
