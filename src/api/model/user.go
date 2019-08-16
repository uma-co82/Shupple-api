package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UID      string    `json:"uid"`
	NickName string    `json:"nickName"`
	Sex      string    `json:"sex"`
	Hobby    string    `json:"hobby"`
	BirthDay time.Time `json:"birthDay"`
}
