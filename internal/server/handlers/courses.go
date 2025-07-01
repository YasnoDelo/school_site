package handlers

import (
	"html/template"
	"net/http"
)

// CoursesData — пример структуры, которую можно передать в шаблон
type CoursesData struct {
	Title   string
	Courses []string
}

func Courses(w http.ResponseWriter, r *http.Request) {
	// Здесь можно собрать реальные данные — из БД или конфига
	data := CoursesData{
		Title:   "Our Courses",
		Courses: []string{"Go Basics", "Advanced Go", "Web Development"},
	}

	// Парсим base + собственный шаблон
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/courses.html",
	))

	// Рендерим базовый шаблон, в него подтянется {{template "content" .}}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
