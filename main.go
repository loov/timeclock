package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
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

	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.HandleFunc("/worker", worker)
	http.HandleFunc("/worker/report", report)
	http.HandleFunc("/worker/working", working)
	http.HandleFunc("/worker/review", review)
	http.HandleFunc("/accountant", accountant)
	http.HandleFunc("/projects", projects)
	http.HandleFunc("/project/", project)

	log.Println("Starting server on", *addr)
	http.ListenAndServe(*addr, nil)
}

type Working struct {
	Activity string
	Started  time.Time
}

func worker(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("start-work.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
	}
}

func working(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("working.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
	}
}

func review(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("review.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
	}
}

func accountant(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("accountant.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
	}
}

func report(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("report.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
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

func projects(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("projects.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}
}

func project(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("project.html", "common.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}
}
