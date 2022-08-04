package dbops

import (
	"testing"
)

func TestAddUser(t *testing.T) {
	err := AddUserCredential("test", "test")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	id, err := GetUserCredential("test", "test")
	if id != 1 || err != nil {
		t.Errorf("Error of GetUser")
	}
	t.Logf("user id: %d\n",id)
}

func TestAddVideoInfo(t *testing.T) {
	err := AddNewVideo("test", "my-video","http://")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
}

func TestGetVideoInfo(t *testing.T) {
	res, err := GetVideoInfo(1)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
	t.Logf("video id: %d\n",res.Video_id)
}

func TestDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(1, "test")
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func TestAddComments(t *testing.T) {
	err := AddNewComments(1, "test", "test")
	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

func TestListComments(t *testing.T) {
	res, err := ListComments(1)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}
	for i, ele := range res {
		t.Logf("comment: %d, %v \n", i, ele)
	}
}

func TestDeleteComments(t *testing.T) {
	err := DeleteCommentInfo(1, "test")
	if err != nil {
		t.Errorf("Error of DeleteComments: %v", err)
	}
}
