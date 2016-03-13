package dayreport

import (
	"net/http"

	"github.com/loov/timeclock/tracking"
)

type Reports interface {
}

type Presenter interface {
	InternalError(w http.ResponseWriter, r *http.Request, err error)
	Present(w http.ResponseWriter, r *http.Request, name string, data interface{})
}

type Server struct {
	Presenter  Presenter
	Activities tracking.Activities
	Reports    Reports
}

func NewServer(presenter Presenter, activities tracking.Activities, reports Reports) *Server {
	return &Server{
		Presenter:  presenter,
		Activities: activities,
		Reports:    reports,
	}
}

func (server *Server) ServeSubmit(w http.ResponseWriter, r *http.Request) {
	//TODO: calculate the totals
	server.Presenter.Present(w, r, "dayreport/submit.html",
		map[string]interface{}{})
}

func (server *Server) ServeList(w http.ResponseWriter, r *http.Request) {
	//TODO: calculate the totals
	server.Presenter.Present(w, r, "dayreport/list.html",
		map[string]interface{}{})
}
