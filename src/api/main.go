package main

import (
	"fmt"
	"github.com/bamzi/jobrunner"
	"github.com/uma-co82/Shupple-api/src/api/db"
	"github.com/uma-co82/Shupple-api/src/api/server"
	"github.com/uma-co82/Shupple-api/src/api/task"
)

func main() {
	db.Init()
	db.AutoMigration()
	server.Init()

	jobrunner.Start()
	jobrunner.Schedule("@every 5s", Myjob{})
}

type Myjob struct {
}

func (e Myjob) Run() {
	fmt.Println("hogehoge")
	task.UserCombinationCheckCreatedAtTask()
}
