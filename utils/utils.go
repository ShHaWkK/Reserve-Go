package utils

import "net/http"
import "html/template"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, "home.html", nil)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	ExecuteTemplate(w, "download.html", nil)
}

func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}
