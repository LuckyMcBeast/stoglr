package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status string

const (
	DISABLED  Status = "DISABLED"
	ENABLED   Status = "ENABLED"
	NOT_FOUND Status = "NOT_FOUND"
	INVALID   Status = "INVALID"
)

type ToggleType string

const (
	RELEASE ToggleType = "RELEASE"
	OPS     ToggleType = "OPS"
	AB      ToggleType = "AB"
)

type Toggle struct {
	Name       string     `json:"name"`
	Status     Status     `json:"status"`
	ToggleType ToggleType `json:"toggleType"`
}

func NewToggle(name string) Toggle {
	return Toggle{Name: name, Status: DISABLED, ToggleType: RELEASE}
}

func CreateOrGetToggle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(NewToggle(name))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(resp))
	w.Write(resp)
}

func createRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/{name}", CreateOrGetToggle)
	return router
}

func main() {

	//toggleDatastore := make(map[string]Toggle)

	s := http.Server{
		Addr:    ":8080",
		Handler: createRouter(),
	}

	log.Fatal(s.ListenAndServe())
}
