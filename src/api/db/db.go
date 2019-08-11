package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func DBConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "shupple"
	PASS := "shupple"
	PROTOCOL := "tcp(db:3314)"
	DBNAME := "shupple"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
	return db
}
