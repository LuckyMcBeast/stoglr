package server

import (
	"io"
	"net/http"
	"stoglr/server/datastore"
	"sync"
	"testing"
)

func TestToggleServer_Start(t *testing.T) {
	port := "3333"
	db := datastore.NewRuntimeDatastore()
	ts := *NewToggleServer(port, db)
	var wg sync.WaitGroup
	wg.Add(1)
	go serverRunner(&ts, &wg)

	client := http.DefaultClient
	resp, err := client.Get("http://localhost:3333/api/health")
	if err != nil {
		t.Fatal("Cannot connect to server", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected: %v, actual: %v", http.StatusOK, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Cannot read body", err)
	}
	if string(b) != "OK" {
		t.Errorf("expected: %v, actual: %v", "OK", string(b))
	}

	wg.Done()
}

func serverRunner(ts *ToggleServer, wg *sync.WaitGroup) {
	defer wg.Done()
	ts.Start()
}
