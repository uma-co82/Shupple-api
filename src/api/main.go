package main

import (
	"./db"
	"./model"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	model.DBMigrate(db.DBConnect())
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ting",
		})
	})
	r.Run(":8080")
}
