package server

import (
	"stoglr/server/datastore"
	"testing"
)

func TestToggleServer_Start(t *testing.T) {
	port := "3333"
	db := datastore.NewRuntimeDatastore()
	ts := NewToggleServer(port, db)

	go ts.Start()
}
