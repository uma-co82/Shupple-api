package db

import (
	"../structs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	err error
)

func Init() {
	DBMS := "mysql"
	USER := "root"
	PASS := "shupple"
	PROTOCOL := "tcp(mysql:3306)"
	DBNAME := "shupple"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
	db, err = gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
	autoMigration()
}

func GetDB() *gorm.DB {
	return db
}

func autoMigration() {
	db.AutoMigrate(&structs.User{})
	db.AutoMigrate(&structs.UserInformation{})
	db.AutoMigrate(&structs.UserCombination{})
	db.AutoMigrate(&structs.InfoCompatible{})
}
