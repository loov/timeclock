package project

import (
	"net/http"

	"github.com/loov/timeclock/site"
)

type Server struct {
	Templates *site.Templates
	Projects  Database
}

func NewServer(templates *site.Templates, projects Database) *Server {
	server := &Server{}
	server.Templates = templates
	server.Projects = projects
	return server
}

func (server *Server) ServeInfos(w http.ResponseWriter, r *http.Request) {
	projects, err := server.Projects.Infos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.Templates.Present(w, r, "project/list.html", map[string]interface{}{
		"Projects": projects,
	})
}
