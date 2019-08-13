package main

import (
	"github.com/gin-gonic/gin"
	"github.com/holefillingco-ltd/Shupple-api/src/api/db"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ting",
		})
	})
	db.Init()
	r.Run(":8080")
}
