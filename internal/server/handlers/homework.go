// internal/server/handlers/homework.go
package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/YasnoDelo/school_site/internal/server/util"
)

type rawTask struct {
	ID       int    `json:"id"`
	Question string `json:"question"` // LaTeX-код
	Answer   string `json:"answer"`
}

type Task struct {
	ID       int
	Question template.HTML
	Answer   string
}

type HomeworkData struct {
	Tasks   []Task
	Results map[int]bool
}

func Homework(w http.ResponseWriter, r *http.Request) {
	// 1) Открываем JSON с задачами
	projectRoot := util.FindProjectRoot()
	dataPath := filepath.Join(projectRoot, "data", "homework.json")
	file, err := os.Open(dataPath)
	if err != nil {
		http.Error(w, "Cannot open tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 2) Парсим JSON
	var raw []rawTask
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		http.Error(w, "Invalid tasks format: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3) Рендерим LaTeX → SVG
	tasks := make([]Task, len(raw))
	for i, t := range raw {
		// оборачиваем LaTeX в блочный MathJax‑делимитер
		tex := fmt.Sprintf("$$%s$$", t.Question)
		tasks[i] = Task{
			ID:       t.ID,
			Question: template.HTML(tex),
			Answer:   t.Answer,
		}
	}

	// 4) Обработка POST‑ответов
	data := HomeworkData{Tasks: tasks}
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Cannot parse form: "+err.Error(), http.StatusBadRequest)
			return
		}
		data.Results = make(map[int]bool)
		for _, t := range tasks {
			userAns := r.FormValue(fmt.Sprintf("answer_%d", t.ID))
			data.Results[t.ID] = (userAns == t.Answer)
		}
	}

	render(w, r, "homework", "Домашние задания", data)
}
