package server

import (
	"context"
	"encoding/json"
	"github.com/a-h/templ"
	"log"
	"net/http"
	"os"
	"stoglr/model"
)

type Handler struct {
	pattern     string
	handlerFunc func(http.ResponseWriter, *http.Request)
}

func (tr *ToggleRouter) getHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Failed to write health check response: " + err.Error())
	}
}

// TODO: Update to produce HTML by default
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

func (tr *ToggleRouter) change(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	t := tr.datastore.ChangeToggle(name)
	writeJsonOrHtml(w, r, model.ToggleHtml(&t), t)
}

func (tr *ToggleRouter) executes(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	_ = r.ParseForm()
	strExe := r.Form.Get("executes")
	if strExe == "" {
		strExe = r.PathValue("executes")
	}
	t := tr.datastore.SetExecution(name, strExe)
	writeJsonOrHtml(w, r, model.ToggleHtml(&t), t)
}

func (tr *ToggleRouter) delete(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	t := tr.datastore.DeleteToggle(name)
	writeJsonOrHtml(w, r, model.ToggleHtml(&t), t)
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

func writeJsonOrHtml(w http.ResponseWriter, r *http.Request, c templ.Component, v any) {
	if r.Header.Get("Accept") == "application/json" {
		writeToJson(w, v)
		return
	}
	writeToHtml(w, c)
}

func writeToHtml(w http.ResponseWriter, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.Render(context.Background(), w)
	writeIfErr(w, err)
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
