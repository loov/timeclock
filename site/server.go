package site

import (
	"net/http"
)

type Server struct {
	Templates *Templates
}

func NewServer(templates *Templates) *Server {
	server := &Server{}
	server.Templates = templates
	return server
}

func (server *Server) ServeEmpty(w http.ResponseWriter, r *http.Request) {
	server.Templates.Present(w, r, "site/empty.html", nil)
}
