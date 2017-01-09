package work

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
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
	Templates  Templates
	activities Activities
}

func NewServer(templates Templates) *Server {
	server := &Server{}
	server.Templates = templates
	server.activities = NewProject("1231", "Railing")
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
		return server.activities.Start(nextActivity)
	} else {
		return server.activities.Finish()
	}

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

		if r.Form.Get("select-activity") == "" {
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

	activityNames, err := server.activities.DefaultNames()
	if err != nil {
		log.Println(err)
	}

	activity, err := server.activities.Current()
	if err != nil {
		log.Println(err)
	}

	server.Templates.Present(w, r, "work/work.html", map[string]interface{}{
		"PostError":    postError.Value,
		"RequestToken": requestToken,

		"CurrentActivity": activity,
		"ActivityNames":   activityNames,
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

	if err := server.activities.Finish(); err != nil {
		return err
	}

	pending, err := server.activities.Pending()
	if err != nil {
		return err
	}

	summary, err := SummarizeActivities(pending)
	if err != nil {
		return err
	}

	err = server.activities.Report(summary)
	if err != nil {
		return err
	}

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

	pending, err := server.activities.Pending()
	if err != nil {
		log.Println(err)
	}

	summary, err := SummarizeActivities(pending)
	if err != nil {
		log.Println(err)
	}

	server.Templates.Present(w, r, "work/submit-report.html", map[string]interface{}{
		"PostError":    postError.Value,
		"RequestToken": requestToken,

		"Pending": pending,
		"Summary": summary,
	})
}

func (server *Server) ServeHistory(w http.ResponseWriter, r *http.Request) {
	reports, err := server.activities.Reports()
	if err != nil {
		log.Println(err)
	}

	server.Templates.Present(w, r, "work/history.html", map[string]interface{}{
		"Reports": reports,
	})
}
