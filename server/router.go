package server

import (
	"html/template"
	"log"
	"net/http"
	"stoglr/model"
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
	return tr.setupUi().
		handleAll(
			Handler{"GET /api/health", tr.getHealth},
			Handler{"GET /api/toggle", tr.getAll},
			Handler{"POST /api/toggle/{name}", tr.createOrGet},
			Handler{"PUT /api/toggle/{name}/change", tr.change},
			Handler{"DELETE /api/toggle/{name}", tr.delete},
			Handler{"PUT /api/toggle/{name}/execute", tr.executes},
			Handler{"PUT /api/toggle/{name}/execute/{executes}", tr.executes},
		).getServeMux()
}

func (tr *ToggleRouter) setupUi() Mux {
	tr.mux.handleAll(
		Handler{
			"GET /",
			func(w http.ResponseWriter, r *http.Request) {
				tmpl := template.Must(template.ParseFiles("static/index.html"))
				err := tmpl.Execute(w, model.NewTogglesByType(tr.datastore.GetAllToggles()))
				if err != nil {
					log.Println(err)
				}
			},
		},
		Handler{
			"GET /simple.min.css",
			func(w http.ResponseWriter, r *http.Request) {
				writeFile("static/simple.min.css", "text/css; charset=utf-8", w)
			},
		},
		Handler{
			"GET /style.css",
			func(w http.ResponseWriter, r *http.Request) {
				writeFile("static/style.css", "text/css; charset=utf-8", w)
			},
		},
		Handler{
			"GET /htmx.min.js",
			func(w http.ResponseWriter, r *http.Request) {
				writeFile("static/htmx.min.js", "text/javascript; charset=utf-8", w)
			},
		},
	)
	return tr.mux
}
