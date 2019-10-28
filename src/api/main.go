package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uma-co82/Shupple-api/src/api/db"
	"github.com/uma-co82/Shupple-api/src/api/server"
)

func main() {
	db.Init()
	db.AutoMigration()
	server.Init()
	gin.SetMode(gin.ReleaseMode)

	//scheduler.Every(1).Minutes().Run(task.UserCombinationCheckCreatedAtTask)
}
