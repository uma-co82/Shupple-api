package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UID      string `json:"uid"`
	NickName string `json:"nickName"`
	Sex      string `json:"sex"`
	Hobby    string `json:"hobby"`
}
