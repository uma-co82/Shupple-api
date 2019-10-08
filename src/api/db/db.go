/***********************************************************
 *                        local                            *
 ***********************************************************/
package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
	"github.com/uma-co82/Shupple-api/src/api/structs"
)

var (
	db  *gorm.DB
	err error
)

type DBENV struct {
	SHUPPLEDBMS      string
	SHUPPLEDBUSER    string
	SHUPPLEDBPASS    string
	SHUPPLEDBPRTOCOL string
	SHUPPLEDBNAME    string
}

func Init() *gorm.DB {
	var env DBENV
	_ = envconfig.Process("", &env)
	DBMS := env.SHUPPLEDBMS
	USER := env.SHUPPLEDBUSER
	PASS := env.SHUPPLEDBPASS
	PROTOCOL := "tcp(" + env.SHUPPLEDBPRTOCOL + ":3306)"
	DBNAME := env.SHUPPLEDBNAME
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

/***********************************************************
 *                    Production                           *
 ***********************************************************/
//package db
//
//import (
////"../structs"
//"github.com/jinzhu/gorm"
//_ "github.com/jinzhu/gorm/dialects/mysql"
//"github.com/uma-co82/Shupple-api/src/api/structs"
//)
//
//var (
//	db  *gorm.DB
//	err error
//)
//
//func Init() *gorm.DB {
//	DBMS := "mysql"
//	USER := "shupple"
//	PASS := "shupple1995"
//	PROTOCOL := "tcp(shupple-api-db.cniwd3cd12wv.ap-northeast-1.rds.amazonaws.com:3306)"
//	DBNAME := "shupple"
//	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
//	db, err = gorm.Open(DBMS, CONNECT)
//	if err != nil {
//		panic(err.Error())
//	}
//	return db
//}
//
//func AutoMigration() {
//	db.AutoMigrate(&structs.User{})
//	db.AutoMigrate(&structs.UserInformation{})
//	db.AutoMigrate(&structs.UserCombination{})
//	db.AutoMigrate(&structs.InfoCompatible{})
//}
