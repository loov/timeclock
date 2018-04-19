package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/loov/timeclock/pgdb"
	"github.com/loov/timeclock/site"
	"github.com/loov/timeclock/user"
	"github.com/loov/timeclock/work"
)

var (
	addr = flag.String("listen", "127.0.0.1:8080", "http server `address`")

	db = flag.String("db", "user=timeclock password=timeclock dbname=timeclock sslmode=disable", "database params")
)

func main() {
	flag.Parse()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host != "" || port != "" {
		*addr = host + ":" + port
	}

	db, err := pgdb.New(*db)
	if err != nil {
		log.Fatal(err)
	}
	err = db.DANGEROUS_DROP_ALL_TABLES()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = db.FakeDatabase()
	if err != nil {
		log.Fatal(err)
	}

	Templates := site.NewTemplates()

	Site := site.NewServer(Templates)
	User := user.NewServer(Templates, db.Users())
	Work := work.NewServer(Templates)

	assets := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.HandleFunc("/", Site.ServeEmpty)
	http.HandleFunc("/user", User.ServeList)
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
