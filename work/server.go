package work

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/site"
	"github.com/loov/timeclock/user"
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
	Templates *site.Templates
	Database  *Database
}

func NewServer(templates *site.Templates) *Server {
	server := &Server{}
	server.Templates = templates
	server.Database = NewDatabase()
	return server
}

func (server *Server) ServeOverview(w http.ResponseWriter, r *http.Request) {
	postError, err := r.Cookie("post-error")
	if err != nil {
		postError = &http.Cookie{}
	}

	http.SetCookie(w, &http.Cookie{Name: "post-error", MaxAge: -1})

	now := time.Now()
	_, _, day := time.Now().Date()

	var start, end time.Time
	if day < 15 {
		start = now.AddDate(0, -1, -day).Truncate(24 * time.Hour)
		end = now.AddDate(0, 1, -day).Truncate(24 * time.Hour)
	} else {
		start = now.AddDate(0, 0, -day).Truncate(24 * time.Hour)
		end = now.AddDate(0, 2, -day).Truncate(24 * time.Hour)
	}

	sheet, err := server.Database.WorkerSheet(1, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.Templates.Present(w, r, "work/overview.html", map[string]interface{}{
		"PostError": postError.Value,
		"Recent":    sheet,
	})
}

var rxActivityField = regexp.MustCompile(`^Activities\[(\d+)\]\.([[:alnum:]]+)$`)

func (server *Server) ServeDaySheet(w http.ResponseWriter, r *http.Request) {
	postError, err := r.Cookie("post-error")
	if err != nil {
		postError = &http.Cookie{}
	}

	const DefaultNumberOfActivitites = 10
	const MaxNumberOfActivities = 50

	activities := make([]Activity, 0, DefaultNumberOfActivitites)

	if r.Method == http.MethodPost {
		r.ParseForm()
		for key, values := range r.Form {
			if len(values) != 1 || (values[0] == "") {
				continue
			}
			value := values[0]

			matches := rxActivityField.FindStringSubmatch(key)
			if len(matches) == 0 {
				continue
			}

			index, err := strconv.Atoi(matches[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if index < 0 || index > MaxNumberOfActivities {
				http.Error(w, "Too many activities.", http.StatusBadRequest)
				return
			}

			if len(activities) <= index {
				activities = append(activities, make([]Activity, index+1-len(activities))...)
			}
			activity := &activities[index]
			activity.Time = time.Now()
			activity.Modified = time.Now()
			activity.Worker = 1

			switch matches[2] {
			case "Project":
				id, _ := strconv.Atoi(value)
				activity.Project = project.ID(id)
			case "Amount":
				amount, _ := strconv.Atoi(value)
				activity.Duration = time.Duration(amount * int(time.Hour))
			case "Activity":
				activity.Name = ActivityName(value)
			default:
				log.Printf("Unknown property %q=%q", key, value)
			}
		}

		server.Database.Submit(activities)
	}

	if len(activities) < DefaultNumberOfActivitites {
		activities = append(activities, make([]Activity, DefaultNumberOfActivitites-len(activities))...)
	}

	http.SetCookie(w, &http.Cookie{Name: "post-error", MaxAge: -1})

	requestToken := createToken()
	http.SetCookie(w, &http.Cookie{
		Path:   "/",
		Name:   "request-token",
		Value:  requestToken,
		MaxAge: 0,
	})

	defaultActivities, err := server.Database.DefaultActivities()
	if err != nil {
		log.Println(err)
	}

	server.Templates.Present(w, r, "work/day-sheet.html", map[string]interface{}{
		"PostError":    postError.Value,
		"RequestToken": requestToken,

		"DefaultActivities": defaultActivities,

		"Worker": user.User{
			ID:   1,
			Name: "Egon Elbre",
		},
		"Projects": []project.Project{
			{ID: 1, Name: "Alpha"},
			{ID: 2, Name: "Beta"},
			{ID: 3, Name: "Gamma"},
			{ID: 4, Name: "Delta"},
		},
		"Activities": activities,
	})
}

/*
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

		http.Redirect(w, r, r.RequestURI+"/..", http.StatusSeeOther)
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
*/
