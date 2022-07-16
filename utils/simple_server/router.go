package simple_server

import (
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{routes: map[string]map[string]HandlerFunc{}}
}

func (r *Router) match(req *http.Request) (HandlerFunc, error) {
	m, ok := r.routes[req.Method]
	if !ok {
		return nil, errors.New("not found")
	}
	uri := strings.Split(req.RequestURI, "?")[0]
	f, ok := m[uri]
	if !ok {
		return nil, errors.New("not found")
	}
	return f, nil
}

func (r *Router) GET(uri string, handlerFunc HandlerFunc) {
	if _, ok := r.routes[http.MethodGet]; !ok {
		r.routes[http.MethodGet] = map[string]HandlerFunc{}
	}
	r.routes[http.MethodGet][uri] = handlerFunc
}

func (r *Router) POST(uri string, handlerFunc HandlerFunc) {
	if _, ok := r.routes[http.MethodPost]; !ok {
		r.routes[http.MethodPost] = map[string]HandlerFunc{}
	}
	r.routes[http.MethodPost][uri] = handlerFunc
}
