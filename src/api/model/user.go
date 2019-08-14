package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	NickName string `json:"nickName"`
}
