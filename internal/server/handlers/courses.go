package handlers

import (
	"net/http"
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

func Courses1(w http.ResponseWriter, r *http.Request) {
	data := CoursesPageData{
		Title: "Наши курсы",
		Courses: []Course{
			{ID: 1, Name: "Math", URL: "/courses/1"},
			{ID: 2, Name: "Physics", URL: "/courses/2"},
		},
	}

	render(w, r, "courses1", "Курсы1", data)
}

func Courses2(w http.ResponseWriter, r *http.Request) {
	data := CoursesPageData{
		Title: "Наши курсы",
		Courses: []Course{
			{ID: 1, Name: "Math", URL: "/courses/1"},
			{ID: 2, Name: "Physics", URL: "/courses/2"},
		},
	}

	render(w, r, "courses2", "Курсы2", data)
}

func Courses3(w http.ResponseWriter, r *http.Request) {
	data := CoursesPageData{
		Title: "Наши курсы",
		Courses: []Course{
			{ID: 1, Name: "Math", URL: "/courses/1"},
			{ID: 2, Name: "Physics", URL: "/courses/2"},
		},
	}

	render(w, r, "courses3", "Курсы3", data)
}
