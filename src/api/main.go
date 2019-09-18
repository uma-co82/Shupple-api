package main

import (
	"./db"
	"./server"
)

func main() {
	db.Init()
	db.AutoMigration()
	server.Init()
}
