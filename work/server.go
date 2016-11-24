package work

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

func createToken() string {
	var data [8]byte
	_, err := rand.Read(data[:])
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(data[:])
}

type Templates interface {
	InternalError(w http.ResponseWriter, r *http.Request, err error)
	Present(w http.ResponseWriter, r *http.Request, name string, data interface{})
}

type Server struct {
	Templates Templates
	project   *Project
}

func NewServer(templates Templates) *Server {
	server := &Server{}
	server.Templates = templates
	server.project = NewProject()
	return server
}

func (server *Server) handleSelectActivity(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	tokenCookie, err := r.Cookie("request-token")
	if err != nil {
		// in case no-cookie, assume it's empty
		tokenCookie = &http.Cookie{}
	}

	tokenForm := r.Form.Get("request-token")
	if tokenForm != tokenCookie.Value && tokenCookie.Value != "" {
		// don't handle refresh
		return nil
	}

	nextActivity := r.Form.Get("select-activity")
	// TODO: validate next activity value
	server.project.SelectActivity(nextActivity)

	return nil
}

func (server *Server) ServeSelectActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := server.handleSelectActivity(w, r)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Path:   "/",
				Name:   "post-error",
				Value:  err.Error(),
				MaxAge: 0,
			})
		}

		if server.project.CurrentActivity() == "" {
			http.Redirect(w, r, r.RequestURI+"/submit", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
		}
		return
	}

	postError, err := r.Cookie("post-error")
	if err != nil {
		postError = &http.Cookie{}
	}

	http.SetCookie(w, &http.Cookie{Name: "post-error", MaxAge: -1})

	requestToken := createToken()
	http.SetCookie(w, &http.Cookie{
		Path:   "/",
		Name:   "request-token",
		Value:  requestToken,
		MaxAge: 0,
	})

	type Data struct {
		PostError    string
		RequestToken string

		CurrentActivity string
		Activities      []string
		Jobs            []Job

		JobSummary map[string]time.Duration
	}

	server.Templates.Present(w, r, "work/select-activity.html", &Data{
		PostError:    postError.Value,
		RequestToken: requestToken,

		CurrentActivity: server.project.CurrentActivity(),
		Activities:      server.project.Activities(),
	})
}

func (server *Server) handleSubmitDay(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	tokenCookie, err := r.Cookie("request-token")
	if err != nil {
		// in case no-cookie, assume it's empty
		tokenCookie = &http.Cookie{}
	}

	tokenForm := r.Form.Get("request-token")
	if tokenForm != tokenCookie.Value && tokenCookie.Value != "" {
		// don't handle refresh
		return nil
	}

	server.project.SubmitDay()

	return nil
}

func (server *Server) ServeSubmitDay(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := server.handleSubmitDay(w, r)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Path:   "/",
				Name:   "post-error",
				Value:  err.Error(),
				MaxAge: 0,
			})
		}

		http.Redirect(w, r, r.RequestURI+"/../history", http.StatusSeeOther)
		return
	}

	postError, err := r.Cookie("post-error")
	if err != nil {
		postError = &http.Cookie{}
	}
	http.SetCookie(w, &http.Cookie{Name: "post-error", MaxAge: -1})

	requestToken := createToken()
	http.SetCookie(w, &http.Cookie{
		Path:   "/",
		Name:   "request-token",
		Value:  requestToken,
		MaxAge: 0,
	})

	type Data struct {
		PostError    string
		RequestToken string

		Jobs       []Job
		JobSummary map[string]time.Duration
	}

	server.Templates.Present(w, r, "work/submit-day.html", &Data{
		PostError:    postError.Value,
		RequestToken: requestToken,

		Jobs:       server.project.Jobs(),
		JobSummary: server.project.Summary(),
	})
}

func (server *Server) ServeHistory(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Days []Day
	}

	server.Templates.Present(w, r, "work/history.html", &Data{
		Days: server.project.Days(),
	})
}
