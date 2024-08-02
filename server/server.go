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
	return &ToggleServer{Port: port, Datastore: datastore}
}

func (t *ToggleServer) Start() {
	log.Println("STOGLR: The Simple Feature Toggler")
	if t.router == nil {
		t.router = NewToggleRouter(t.Datastore)
	}
	if t.server == nil {
		t.server = &http.Server{
			Addr:    ":" + t.Port,
			Handler: t.router.CreateRouter(),
		}
	}
	log.Println("Starting STOGLR on port", t.Port)
	log.Fatal(t.server.ListenAndServe())
}
