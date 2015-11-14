package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/loov/timeclock/project"
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

	http.HandleFunc("/", index)

	log.Println("Starting server on", *addr)
	http.ListenAndServe(*addr, nil)
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

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}

	example := &project.Project{
		Title:    "Alpha",
		Customer: "ACME",
		Pricing: project.Pricing{
			Hours: 480,
			Price: 1000,
		},
		Description: "Implement views",
		Status:      project.InProgress,
	}

	err = t.Execute(w, example)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		internalError(w, r, err)
		return
	}
}
