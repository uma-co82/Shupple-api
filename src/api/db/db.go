package db

import (
	//"../structs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/uma-co82/Shupple-api/src/api/structs"
)

var (
	db  *gorm.DB
	err error
)

func Init() *gorm.DB {
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
	return db
}

func AutoMigration() {
	db.AutoMigrate(&structs.User{})
	db.AutoMigrate(&structs.UserInformation{})
	db.AutoMigrate(&structs.UserCombination{})
	db.AutoMigrate(&structs.InfoCompatible{})
}
