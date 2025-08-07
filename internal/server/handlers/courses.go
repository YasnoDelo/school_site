// internal/server/handlers/courses.go
package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Topic описывает одну тему внутри курса
type Topic struct {
	Name      string // например "Algebra"
	Methodics string // имя файла методички, например "algebra.pdf"
	Lecture   string // имя файла лекции, e.g. "algebra.mp4"
	Cheat     string // имя шпаргалки, e.g. "algebra.png"
}

// CourseLink — для списка всех курсов
type CourseLink struct {
	ID      string // "math", "infa" и т.д.
	Title   string // "Математика"
	Summary string // краткое описание
}

// ListCoursesData — payload для страницы /courses
type ListCoursesData struct {
	Courses []CourseLink
}

// CoursePageData — payload для /courses/{subject}
type CoursePageData struct {
	Subject string  // тот же ID, e.g. "math"
	Title   string  // читаемое название, e.g. "Математика"
	Topics  []Topic // список тем
}

// ListCourses отдаёт страницу со списком курсов
func ListCourses(w http.ResponseWriter, r *http.Request) {
	data := ListCoursesData{
		Courses: []CourseLink{
			{ID: "math", Title: "Математика", Summary: "Алгебра, геометрия и т.д."},
			{ID: "infa", Title: "Информатика", Summary: "Программирование, алгоритмы"},
		},
	}
	render(w, r,
		"courses_list", // <-- имя шаблона courses_list.html
		"Курсы",        // <-- заголовок страницы
		data,           // <-- ваш payload
	)
}

// CoursePage отдаёт детальную страницу конкретного курса
func CoursePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"]

	var title string
	var topics []Topic

	switch subject {
	case "math":
		title = "Математика"
		topics = []Topic{
			{Name: "Intro", Methodics: "intro.pdf", Lecture: "intro.mp4", Cheat: "intro.png"},
			{Name: "Algebra", Methodics: "algebra.pdf", Lecture: "algebra.mp4", Cheat: "algebra.png"},
		}
	case "infa":
		title = "Информатика"
		topics = []Topic{
			{Name: "Basics", Methodics: "basics.pdf", Lecture: "basics.mp4", Cheat: "basics.png"},
			{Name: "Algorithms", Methodics: "algo.pdf", Lecture: "algo.mp4", Cheat: "algo.png"},
		}
	default:
		http.NotFound(w, r)
		return
	}

	render(w, r,
		"courses_detail", // имя шаблона courses_detail.html
		title,            // заголовок страницы, например "Математика"
		CoursePageData{
			Subject: subject,
			Title:   title,
			Topics:  topics,
		},
	)
}
