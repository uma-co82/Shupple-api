package main

import (
	"github.com/uma-co82/Shupple-api/src/api/db"
	"github.com/uma-co82/Shupple-api/src/api/server"
)

func main() {
	db.Init()
	db.AutoMigration()
	server.Init()
}
