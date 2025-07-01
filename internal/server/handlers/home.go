package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func Home(w http.ResponseWriter, r *http.Request) {
	// Собираем полный путь к файлам
	basePath := filepath.Join(TemplatesDir, "base.html")
	homePath := filepath.Join(TemplatesDir, "home.html")

	tmpl := template.Must(template.ParseFiles(basePath, homePath))

	// первый аргумент — имя шаблона из base.html (<html>{{template "content"}}</html>)
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
	}
}
