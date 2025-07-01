package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// MaterialsData — данные для материалов
type MaterialsData struct {
	Title     string
	FileNames []string
}

func Materials(w http.ResponseWriter, r *http.Request) {
	// Предположим, все файлы лежат в папке static/materials/
	files, err := filepath.Glob("static/materials/*")
	if err != nil {
		http.Error(w, "Failed to list materials", http.StatusInternalServerError)
		return
	}

	data := MaterialsData{
		Title:     "Materials",
		FileNames: files, // сюда попадут пути вроде "static/materials/lesson1.pdf"
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/materials.html",
	))

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
