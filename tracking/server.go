package tracking

import (
	"log"
	"net/http"
	"path"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/user"
)

type Tracker interface {
	SelectActivity(user user.ID, project project.ID, activity string) error
	FinishActivity(user user.ID, project project.ID, activity string) error
}

type Activities interface {
	Unreported() ([]Activity, error)
	UnreportedByUser(user user.ID) ([]Activity, error)

	MarkProcessed(activity []Activity) error
}

type Presenter interface {
	InternalError(w http.ResponseWriter, r *http.Request, err error)
	Present(w http.ResponseWriter, r *http.Request, name string, data interface{})
}

type Server struct {
	Presenter  Presenter
	Tracker    Tracker
	Activities Activities
	Projects   project.Projects
}

func NewServer(presenter Presenter, tracker Tracker, activities Activities, projects project.Projects) *Server {
	return &Server{
		Presenter:  presenter,
		Tracker:    tracker,
		Activities: activities,
		Projects:   projects,
	}
}

func (server *Server) ServeSelectProject(w http.ResponseWriter, r *http.Request) {
	// only active projects
	items, err := server.Projects.List()
	if err != nil {
		log.Printf("error loading projects list: %v", err)
		server.Presenter.InternalError(w, r, err)
		return
	}

	// get actual logged in user
	unreported, err := server.Activities.UnreportedByUser(0)
	if err != nil {
		log.Printf("error loading unprocessed activities: %v", err)
		server.Presenter.InternalError(w, r, err)
		return
	}

	server.Presenter.Present(w, r, "tracking/select-project.html",
		map[string]interface{}{
			"Projects":             items,
			"UnreportedActivities": unreported,
		})
}

func (server *Server) ServeActiveProject(w http.ResponseWriter, r *http.Request) {
	// TODO: form value to the last/current active
	server.ServeSelectActivity(w, r)
}

func (server *Server) ServeSelectActivity(w http.ResponseWriter, r *http.Request) {
	id := project.ID(r.FormValue("id"))
	if id == "" {
		id = project.ID(path.Base(r.URL.Path))
	}

	proj, err := server.Projects.ByID(id)
	if err == project.ErrNotExist {
		//TODO: redirect to select-project.html with flash error
		server.Presenter.InternalError(w, r, err)
		return
	}
	if err != nil {
		log.Printf("error loading project info: %v", err)
		//TODO: redirect to select-project.html with flash error
		server.Presenter.InternalError(w, r, err)
		return
	}

	//TODO: add active activity and project
	server.Presenter.Present(w, r, "tracking/select-activity.html",
		map[string]interface{}{
			"Project":       proj,
			"ActivityNames": ActivityNames,
		})
}
