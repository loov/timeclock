package site

import (
	"compress/gzip"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Templates struct{}

func NewTemplates() *Templates {
	return &Templates{}
}

func (templates *Templates) InternalError(w http.ResponseWriter, r *http.Request, err error) {
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

func (templates *Templates) Present(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	funcs := template.FuncMap{
		"FormatDay": func(t time.Time) string {
			return t.Format("02 Jan 2006")
		},
		"RequestPath": func() string {
			return r.URL.Path
		},
	}

	t, err := template.New("").Funcs(funcs).ParseFiles("site/common.html", name)
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
