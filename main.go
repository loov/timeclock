package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/loov/timeclock/project"
	"github.com/loov/timeclock/project/projects-sql"
)

var (
	addr = flag.String("listen", ":8000", "http server `address`")
)

func main() {
	flag.Parse()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host != "" || port != "" {
		*addr = host + ":" + port
	}

	ProjectDB, err := projects.New("main.db")
	if err != nil {
		log.Fatal(err)
	}

	Project := project.NewServer(Templates{}, ProjectDB)

	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.HandleFunc("/worker", Template("start-work.html", nil))
	http.HandleFunc("/worker/report", Template("report.html", nil))
	http.HandleFunc("/worker/working", Template("working.html", nil))
	http.HandleFunc("/worker/review", Template("review.html", nil))
	http.HandleFunc("/accountant", Template("accountant.html", nil))

	http.HandleFunc("/projects", Project.ServeList)
	http.HandleFunc("/project/", Project.ServeProjectInfo)
	http.HandleFunc("/", Project.ServeList)

	log.Println("Starting server on", *addr)
	http.ListenAndServe(*addr, nil)
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
	t, err := template.ParseFiles(name, "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		templates.InternalError(w, r, err)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
	}
}

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
