package main

import (
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/db"
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/server"
)

func main() {
	db.Init()
	db.AutoMigration()
	//go func() {
	//	scheduler.Every(1).Minutes().Run(task.UserCombinationCheckCreatedAtTask)
	//}()
	server.Init()
}
