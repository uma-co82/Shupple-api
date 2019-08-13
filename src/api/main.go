package main

import (
	"github.com/holefillingco-ltd/Shupple-api/src/api/db"
	"github.com/holefillingco-ltd/Shupple-api/src/api/server"
)

func main() {
	db.Init()
	server.Init()
}
