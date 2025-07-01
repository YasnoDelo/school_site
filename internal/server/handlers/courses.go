package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func Courses(w http.ResponseWriter, r *http.Request) {
	basePath := filepath.Join(TemplatesDir, "base.html")
	coursesPath := filepath.Join(TemplatesDir, "courses.html")

	tmpl := template.Must(template.ParseFiles(basePath, coursesPath))
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
	}
}
