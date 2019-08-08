package main

import (
	"./db"
	"./model"
	"github.com/gin-gonic/gin"
)

func main() {
	model.DBMigrate(db.DBConnect())
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	r.Run(":8080")
}
