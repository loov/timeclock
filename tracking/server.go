package tracking

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"sync"

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

	mu              sync.Mutex
	currentActivity string
}

func NewServer(templates Templates, db *db.DB) *Server {
	server := &Server{}
	server.Templates = templates
	server.DB = db
	server.currentActivity = "Welding"
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
	if nextActivity != "" {
		server.mu.Lock()
		server.currentActivity = nextActivity
		server.mu.Unlock()
	} else {
		// TODO: invalid activity
	}

	return nil
}

func (server *Server) ServeSelectActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := server.handleSelectActivity(w, r)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
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
		Name:   "request-token",
		Value:  requestToken,
		MaxAge: 0,
	})

	type Data struct {
		PostError    string
		RequestToken string

		CurrentActivity string
		Activities      []string
	}

	server.mu.Lock()
	currentActivity := server.currentActivity
	server.mu.Unlock()

	server.Templates.Present(w, r, "tracking/select-activity.html", &Data{
		PostError:    postError.Value,
		RequestToken: requestToken,

		CurrentActivity: currentActivity,
		Activities:      []string{"Plumbing", "Welding", "Construction"},
	})
}
