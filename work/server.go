package work

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/loov/timeclock/db"
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
	DB        *db.DB

	mu   sync.Mutex
	jobs []Job
	days []Day
}

type Job struct {
	Activity string
	Start    time.Time
	Finish   time.Time
}

func (job *Job) Duration() time.Duration {
	if job.Finish.IsZero() {
		return time.Now().Sub(job.Start)
	}
	return job.Finish.Sub(job.Start)
}

type Day struct {
	Submitted  time.Time
	Activities map[string]time.Duration
}

func (server *Server) selectActivity(activity string) {
	server.mu.Lock()
	defer server.mu.Unlock()

	now := time.Now()
	if len(server.jobs) > 0 {
		last := &server.jobs[len(server.jobs)-1]
		if last.Finish.IsZero() {
			last.Finish = now
		}
	}

	if activity != "" {
		server.jobs = append(server.jobs, Job{
			Activity: activity,
			Start:    time.Now(),
		})
	}
}

func (server *Server) submitDay() {
	server.mu.Lock()
	defer server.mu.Unlock()

	durations := map[string]time.Duration{}
	for _, job := range server.jobs {
		durations[job.Activity] += job.Duration()
	}

	day := Day{
		Submitted:  time.Now(),
		Activities: durations,
	}

	server.days = append(server.days, day)
	server.jobs = nil
}

func (server *Server) clonejobs() []Job {
	server.mu.Lock()
	defer server.mu.Unlock()

	return append([]Job{}, server.jobs...)
}

func (server *Server) clonedays() []Day {
	server.mu.Lock()
	defer server.mu.Unlock()

	return append([]Day{}, server.days...)
}

func (server *Server) summarizejobs() map[string]time.Duration {
	server.mu.Lock()
	defer server.mu.Unlock()

	durations := map[string]time.Duration{}
	for _, job := range server.jobs {
		durations[job.Activity] += job.Duration()
	}
	return durations
}

func (server *Server) currentActivity() string {
	server.mu.Lock()
	defer server.mu.Unlock()

	if len(server.jobs) == 0 {
		return ""
	}

	last := &server.jobs[len(server.jobs)-1]
	if last.Finish.IsZero() {
		return last.Activity
	}

	return ""
}

func NewServer(templates Templates, db *db.DB) *Server {
	server := &Server{}
	server.Templates = templates
	server.DB = db
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
	server.selectActivity(nextActivity)

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

		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
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

		CurrentActivity: server.currentActivity(),
		Activities:      []string{"Plumbing", "Welding", "Construction"},
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

	server.submitDay()

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

		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
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

		Jobs:       server.clonejobs(),
		JobSummary: server.summarizejobs(),
	})
}

func (server *Server) ServeHistory(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Days []Day
	}

	server.Templates.Present(w, r, "work/history.html", &Data{
		Days: server.clonedays(),
	})
}
