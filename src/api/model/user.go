package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DBへ保存する為の構造体
type User struct {
	gorm.Model
	UID      string `json:"uid"`
	NickName string `json:"nickName"`
	Sex      string `json:"sex"`
	Hobby    string `json:"hobby"`
	Age      int    `json:"Age"`
}

// POSTされた値を受け取る為の構造体
type PostUser struct {
	NickName string    `json:"nickName"`
	Sex      string    `json:"sex"`
	Hobby    string    `json:"hobby"`
	BirthDay time.Time `json:"birthDay"`
}
