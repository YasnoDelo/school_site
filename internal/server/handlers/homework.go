package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type rawTask struct {
	ID       int    `json:"id"`
	Question string `json:"question"` // LaTeX-код
	Answer   string `json:"answer"`
}

type Task struct {
	ID       int
	Question template.HTML // теперь здесь будем хранить LaTeX
	Answer   string
}

type ViewData struct {
	Tasks   []Task
	Results map[int]bool
}

func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic("cannot get cwd: " + err.Error())
	}
	dir := cwd
	for i := 0; i < 4; i++ {
		if fi, err := os.Stat(filepath.Join(dir, "data")); err == nil && fi.IsDir() {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	panic("project root not found from " + cwd)
}

func Homework(w http.ResponseWriter, r *http.Request) {
	// 1) Открываем JSON с задачами
	projectRoot := findProjectRoot()
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
	data := ViewData{Tasks: tasks}
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

	// 5) Рендерим шаблон
	base := filepath.Join(TemplatesDir, "base.html")
	hw := filepath.Join(TemplatesDir, "homework.html")
	tmpl := template.Must(template.ParseFiles(base, hw))
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
