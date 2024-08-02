package server

import (
	"encoding/json"
	"log"
	"net/http"
	"stoglr/server/datastore"
)

type ToggleRouter struct {
	datastore datastore.Datastore
	mux       Mux
}

func NewToggleRouter(datastore *datastore.RuntimeDatastore) *ToggleRouter {
	return &ToggleRouter{datastore: datastore, mux: NewMuxWrapper()}
}

func (tr *ToggleRouter) CreateRouter() *http.ServeMux {
	tr.mux.handleFunc("GET /api/health", tr.getHealth)
	tr.mux.handleFunc("GET /api/toggle", tr.getAll)
	tr.mux.handleFunc("POST /api/toggle/{name}", tr.createOrGet)
	tr.mux.handleFunc("PUT /api/toggle/{name}/enable", tr.enable)
	tr.mux.handleFunc("PUT /api/toggle/{name}/disable", tr.disable)
	tr.mux.handleFunc("PUT /api/toggle/{name}/execute/{executes}", tr.executes)
	tr.mux.handleFunc("DELETE /api/toggle/{name}", tr.delete)
	return tr.mux.getServeMux()
}

func (tr *ToggleRouter) getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Failed to write health check response: " + err.Error())
	}
}

func (tr *ToggleRouter) getAll(w http.ResponseWriter, r *http.Request) {
	writeToJson(w, tr.datastore.GetAllToggles())
}

func (tr *ToggleRouter) createOrGet(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	queryParams := r.URL.Query()
	toggleType := queryParams.Get("type")
	executes := queryParams.Get("executes")
	writeToJson(w, tr.datastore.CreateOrGetToggle(name, toggleType, executes))
}

func (tr *ToggleRouter) enable(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	writeToJson(w, tr.datastore.EnableToggle(name))
}

func (tr *ToggleRouter) disable(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	writeToJson(w, tr.datastore.DisableToggle(name))
}

func (tr *ToggleRouter) executes(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	strExe := r.PathValue("executes")
	writeToJson(w, tr.datastore.SetExecution(name, strExe))
}

func (tr *ToggleRouter) delete(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	writeToJson(w, tr.datastore.DeleteToggle(name))
}

func writeToJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(v)
	writeIfErr(w, err)
	write(w, resp)
}

func writeIfErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		write(w, []byte(err.Error()))
	}
}

func write(w http.ResponseWriter, resp []byte) {
	_, err := w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
