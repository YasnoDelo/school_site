package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Course struct {
	ID   int
	Name string
	URL  string
}

type CoursesPageData struct {
	Title   string
	Courses []Course
}

func Courses(w http.ResponseWriter, r *http.Request) {
	basePath := filepath.Join(TemplatesDir, "base.html")
	coursesPath := filepath.Join(TemplatesDir, "courses.html")

	data := CoursesPageData{
		Title: "Наши курсы",
		Courses: []Course{
			{ID: 1, Name: "Math", URL: "/courses/1"},
			{ID: 2, Name: "Physics", URL: "/courses/2"},
		},
	}

	tmpl := template.Must(template.ParseFiles(basePath, coursesPath))
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
	}
}
