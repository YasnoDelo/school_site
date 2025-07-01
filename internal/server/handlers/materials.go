package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func Materials(w http.ResponseWriter, r *http.Request) {
	basePath := filepath.Join(TemplatesDir, "base.html")
	materialsPath := filepath.Join(TemplatesDir, "materials.html")

	tmpl := template.Must(template.ParseFiles(basePath, materialsPath))
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
	}
}
