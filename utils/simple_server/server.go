package simple_server

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Server struct {
	router *Router
}

func (h *Server) SetRoute(router *Router) {
	h.router = router
}

func (h *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	f, err := h.router.match(req)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	f(NewContext(w, req))
}
