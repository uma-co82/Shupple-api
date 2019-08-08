package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	NickName string `json:"nickName"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{})
	return db
}
