package project

import (
	"log"
	"net/http"
	"path"
)

type Projects interface {
	List() ([]Project, error)
	ByID(id ID) (Project, error)
}

type Presenter interface {
	InternalError(w http.ResponseWriter, r *http.Request, err error)
	Present(w http.ResponseWriter, r *http.Request, name string, data interface{})
}

type Server struct {
	Presenter Presenter
	Projects  Projects
}

func NewServer(presenter Presenter, projects Projects) *Server {
	return &Server{
		Presenter: presenter,
		Projects:  projects,
	}
}

func (server *Server) ServeList(w http.ResponseWriter, r *http.Request) {
	items, err := server.Projects.List()
	if err != nil {
		log.Printf("error loading projects list: %v", err)
		server.Presenter.InternalError(w, r, err)
		return
	}
	server.Presenter.Present(w, r, "project/list.html", map[string]interface{}{
		"Projects":         items,
		"InactiveProjects": items,
	})
}

func (server *Server) ServeInfo(w http.ResponseWriter, r *http.Request) {
	id := ID(r.FormValue("id"))
	if id == "" {
		id = ID(path.Base(r.URL.Path))
	}

	project, err := server.Projects.ByID(id)
	if err == ErrNotExist {
		//TODO: redirect to select-project.html with flash error
		server.Presenter.InternalError(w, r, err)
		return
	}
	if err != nil {
		log.Printf("error loading project info: %v", err)
		server.Presenter.InternalError(w, r, err)
		return
	}

	server.Presenter.Present(w, r, "project/info.html", map[string]interface{}{
		"Project": project,
	})
}

func (server *Server) ServeAdd(w http.ResponseWriter, r *http.Request) {
	server.Presenter.Present(w, r, "project/add.html", map[string]interface{}{
		"Project": &Project{},
	})
}
