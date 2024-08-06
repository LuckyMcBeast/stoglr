package server

import (
	"net/http"
)

type Mux interface {
	handleFunc(string, func(http.ResponseWriter, *http.Request)) Mux
	getServeMux() *http.ServeMux
	handleAll(...Handler) Mux
}

type MuxWrapper struct {
	serveMux *http.ServeMux
}

func NewMuxWrapper() MuxWrapper {
	return MuxWrapper{serveMux: http.NewServeMux()}
}

func (r MuxWrapper) handleAll(funcs ...Handler) Mux {
	for _, h := range funcs {
		r.serveMux.HandleFunc(h.pattern, h.handlerFunc)
	}
	return r
}

func (r MuxWrapper) handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) Mux {
	r.serveMux.HandleFunc(pattern, handler)
	return r
}

func (r MuxWrapper) getServeMux() *http.ServeMux {
	return r.serveMux
}
