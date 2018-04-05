package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/loov/timeclock/work"
)

var (
	addr = flag.String("listen", "127.0.0.1:8080", "http server `address`")
)

func main() {
	flag.Parse()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host != "" || port != "" {
		*addr = host + ":" + port
	}

	templates := Templates{}

	Work := work.NewServer(templates)

	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.HandleFunc("/work", Work.ServeOverview)
	http.HandleFunc("/work/day", Work.ServeDaySheet)

	// http.HandleFunc("/work/submit", Work.ServeSubmitDay)

	/*
		templates := Templates{}

		DB, err := db.New("main.db")
		if err != nil {
			log.Fatal(err)
		}

		Project := project.NewServer(templates, DB.Projects())
		Tracking := tracking.NewServer(templates, DB.Tracker(), DB.Activities(), DB.Projects())
		DayReport := dayreport.NewServer(templates, DB.Activities(), DB.DayReports())

		http.HandleFunc("/track/project", Tracking.ServeSelectProject)
		http.HandleFunc("/track/active", Tracking.ServeActiveProject)
		http.HandleFunc("/track/project/", Tracking.ServeSelectActivity)

		http.HandleFunc("/projects", Project.ServeList)
		http.HandleFunc("/project/add", Project.ServeAdd)
		http.HandleFunc("/project/", Project.ServeInfo)
		http.HandleFunc("/", Project.ServeList)

		http.HandleFunc("/day/submit", DayReport.ServeSubmit)
		http.HandleFunc("/day/reports", DayReport.ServeList)

		http.HandleFunc("/worker/review", Template("review.html", nil))
		http.HandleFunc("/accountant", Template("accountant.html", nil))
	*/

	http.HandleFunc("/favicon.ico", ServeFavIcon)

	log.Println("Starting server on", *addr)
	log.Println(http.ListenAndServe(*addr, nil))
}

func ServeFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.FromSlash("assets/favicon.png"))
}

type Templates struct{}

func (templates Templates) InternalError(w http.ResponseWriter, r *http.Request, err error) {
	message := template.HTMLEscapeString(err.Error())
	page := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Timeclock</title>
	<link rel="stylesheet" href="/assets/css/main.css">
</head>
<body>
	<div class="error">%s</div>
</body>
</html>
`, message)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(page))
}

func (templates Templates) Present(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	funcs := template.FuncMap{
		"FormatDay": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
		"RequestPath": func() string {
			return r.URL.Path
		},
	}

	t, err := template.New("").Funcs(funcs).ParseFiles("common.html", name)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		templates.InternalError(w, r, err)
		return
	}

	var dest io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err == nil {
			w.Header().Set("Content-Encoding", "gzip")
			defer gz.Close()
			dest = gz
		}
	}

	err = t.ExecuteTemplate(dest, filepath.Base(name), data)
	if err != nil {
		log.Printf("error executing template: %v", err)
	}
}

/*

type Working struct {
	Activity string
	Started  time.Time
}

func Template(name string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(name, "common.html")
		if err != nil {
			log.Printf("error parsing template: %v", err)
			internalError(w, r, err)
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Printf("error executing template: %v", err)
		}
	}
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	message := template.HTMLEscapeString(err.Error())
	page := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Timeclock</title>
	<link rel="stylesheet" href="/assets/css/main.css">
</head>
<body>
	<div class="error">%s</div>
</body>
</html>
`, message)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(page))
}

func ServeFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.FromSlash("assets/favicon.png"))
}
*/
