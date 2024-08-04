package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	tr.setupUi()
	tr.mux.handleFunc("GET /api/health", tr.getHealth)
	tr.mux.handleFunc("GET /api/toggle", tr.getAll)
	tr.mux.handleFunc("POST /api/toggle/{name}", tr.createOrGet)
	tr.mux.handleFunc("PUT /api/toggle/{name}/enable", tr.enable)
	tr.mux.handleFunc("PUT /api/toggle/{name}/disable", tr.disable)
	tr.mux.handleFunc("PUT /api/toggle/{name}/execute/{executes}", tr.executes)
	tr.mux.handleFunc("DELETE /api/toggle/{name}", tr.delete)
	return tr.mux.getServeMux()
}

func (tr *ToggleRouter) setupUi() {
	tr.mux.handleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeFile("static/index.html", "text/html; charset=utf-8", w)
	})
	tr.mux.handleFunc("GET /simple.min.css", func(w http.ResponseWriter, r *http.Request) {
		writeFile("static/simple.min.css", "text/css; charset=utf-8", w)
	})
	tr.mux.handleFunc("GET /htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		writeFile("static/htmx.min.js", "text/javascript; charset=utf-8", w)
	})
}

func (tr *ToggleRouter) getHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Failed to write health check response: " + err.Error())
	}
}

func (tr *ToggleRouter) getAll(w http.ResponseWriter, _ *http.Request) {
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

func writeFile(f string, c string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", c)
	file, err1 := os.Open(f)
	if err1 != nil {
		log.Println("Failed to open file: " + f + ": " + err1.Error())
	}
	stat, err2 := file.Stat()
	if err2 != nil {
		log.Println("Failed to stat file: " + f + ": " + err2.Error())
		panic(err2)
	}
	b := make([]byte, stat.Size())
	_, err3 := file.Read(b)
	if err3 != nil {
		log.Println("Failed to read file: " + f + ": " + err3.Error())
	}
	_, err4 := w.Write(b)
	if err4 != nil {
		log.Println("Failed to write health check response: " + err4.Error())
	}
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
