package main

import (
	"github.com/carlescere/scheduler"
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/db"
	"github.com/uma-co82/Shupple-api/src/api/infrastructure/server"
	"github.com/uma-co82/Shupple-api/src/api/task"
)

func main() {
	db.Init()
	db.AutoMigration()
	go func() {
		scheduler.Every(1).Minutes().Run(task.UserCombinationCheckCreatedAtTask)
	}()
	server.Init()
}
