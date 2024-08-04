package server

import (
	"log"
	"net/http"
	"stoglr/server/datastore"
)

type ToggleServer struct {
	Port      string
	Datastore *datastore.RuntimeDatastore
	server    *http.Server
	router    *ToggleRouter
}

func NewToggleServer(port string, datastore *datastore.RuntimeDatastore) *ToggleServer {
	router := NewToggleRouter(datastore)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router.CreateRouter(),
	}
	return &ToggleServer{Port: port, Datastore: datastore, router: router, server: server}
}

func (t *ToggleServer) Start() {
	log.Println("STOGLR: The Simple Feature Toggler")
	log.Println("Starting STOGLR Server on port", t.Port)
	log.Fatal(t.server.ListenAndServe())
}
