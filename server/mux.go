package server

import (
	"net/http"
)

type Mux interface {
	handleFunc(string, func(http.ResponseWriter, *http.Request))
	getServeMux() *http.ServeMux
}

type MuxWrapper struct {
	serveMux *http.ServeMux
}

func NewMuxWrapper() MuxWrapper {
	return MuxWrapper{serveMux: http.NewServeMux()}
}

func (r MuxWrapper) handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.serveMux.HandleFunc(pattern, handler)
}

func (r MuxWrapper) getServeMux() *http.ServeMux {
	return r.serveMux
}
