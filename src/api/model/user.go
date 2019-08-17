package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DBへ保存する為の構造体
type User struct {
	gorm.Model
	UID      string    `json:"uid"`
	NickName string    `json:"nickName"`
	Sex      int       `json:"sex"`
	Hobby    string    `json:"hobby"`
	BirthDay time.Time `json:"birthDay"`
	Age      int       `json:"Age"`
}

// POSTされた値を受け取る為の構造体
type PostUser struct {
	UID      string    `json:"uid"`
	NickName string    `json:"nickName"`
	Sex      int       `json:"sex"`
	Hobby    string    `json:"hobby"`
	BirthDay time.Time `json:"birthDay"`
}
