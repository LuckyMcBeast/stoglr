package main

import (
	"stoglr/server"
	"stoglr/server/datastore"
)

func main() {
	server.NewToggleServer(
		"8080",
		datastore.NewRuntimeDatastore()).Start()
}
