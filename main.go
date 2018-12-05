package main

import (
	"PentagoServer/server"
	"PentagoServer/db"
)

func main() {
	//init di
	db.Connect()
	//start http server
	server.Start()
}
