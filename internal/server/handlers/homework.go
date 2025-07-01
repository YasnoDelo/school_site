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
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Task struct {
	ID       int
	Question template.HTML
	Answer   string
}

// ViewData передаётся в шаблон: список задач и (опционально) результаты
type ViewData struct {
	Tasks   []Task
	Results map[int]bool
}

// findProjectRoot ищет папку data, поднимаясь вверх
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

	// 2) Парсим JSON в промежуточную структуру rawTask
	var raw []rawTask
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		http.Error(w, "Invalid tasks format: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3) Конвертируем rawTask → Task (с template.HTML для LaTeX)
	tasks := make([]Task, len(raw))
	for i, t := range raw {
		tasks[i] = Task{
			ID:       t.ID,
			Question: template.HTML(t.Question),
			Answer:   t.Answer,
		}
	}

	// 4) Собираем данные для шаблона, обрабатываем POST‑ответы
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
	basePath := filepath.Join(TemplatesDir, "base.html")
	hwPath := filepath.Join(TemplatesDir, "homework.html")
	tmpl := template.Must(template.ParseFiles(basePath, hwPath))
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Template render error: "+err.Error(), http.StatusInternalServerError)
	}
}
