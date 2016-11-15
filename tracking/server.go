package tracking

import (
	"net/http"

	"github.com/loov/timeclock/db"
)

type Templates interface {
	InternalError(w http.ResponseWriter, r *http.Request, err error)
	Present(w http.ResponseWriter, r *http.Request, name string, data interface{})
}

type Server struct {
	Templates Templates
	DB        *db.DB

	CurrentActivity string
}

func NewServer(templates Templates, db *db.DB) *Server {
	return &Server{templates, db, "Welding"}
}

func (server *Server) ServeSelectActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err == nil {
			if nextActivity := r.Form.Get("select-activity"); nextActivity != "" {
				server.CurrentActivity = nextActivity
			}
		} else {
			server.Templates.InternalError(w, r, err)
		}
	}

	type Data struct {
		CurrentActivity string
		Activities      []string
	}

	server.Templates.Present(w, r, "tracking/select-activity.html", &Data{
		CurrentActivity: server.CurrentActivity,
		Activities:      []string{"Plumbing", "Welding", "Construction"},
	})
}
