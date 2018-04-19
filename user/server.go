package user

import (
	"net/http"

	"github.com/loov/timeclock/site"
)

type Server struct {
	Templates *site.Templates
	Users     Database
}

func NewServer(templates *site.Templates, users Database) *Server {
	server := &Server{}
	server.Templates = templates
	server.Users = users
	return server
}

func (server *Server) ServeList(w http.ResponseWriter, r *http.Request) {
	users, err := server.Users.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.Templates.Present(w, r, "user/list.html", map[string]interface{}{
		"Users": users,
	})
}
