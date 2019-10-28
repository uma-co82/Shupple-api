package main

import (
	"github.com/carlescere/scheduler"
	"github.com/uma-co82/Shupple-api/src/api/db"
	"github.com/uma-co82/Shupple-api/src/api/server"
	"github.com/uma-co82/Shupple-api/src/api/task"
)

func main() {
	db.Init()
	db.AutoMigration()
	server.Init()
	scheduler.Every(1).Minutes().Run(task.UserCombinationCheckCreatedAtTask)
}
