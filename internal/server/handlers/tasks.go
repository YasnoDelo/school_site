package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/gorilla/mux"
)

// Problem описывает одну задачу
type Problem struct {
	ID       int      `json:"id"`
	Topic    string   `json:"topic"`
	Prompt   string   `json:"prompt"`
	Images   []string `json:"images"`
	Answers  []string `json:"answers"`
	Solution string   `json:"solution"`
}

// TasksData — данные для шаблона
type TasksData struct {
	Subject string
	Topic   string
	Tasks   []Problem
	Topics  []string
	Checked bool           // ← был ли POST
	Results map[int]bool   // ← результат по каждой задаче: id → правильность
	UserAns map[int]string // ← что ввёл пользователь
}

// tasksByTopic читает банк из static/courses/{subject}/problems/bank.json
func TasksByTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"]
	topic := vars["topic"]

	bankPath := filepath.Join("static", "courses", subject, "problems", "bank.json")
	raw, err := ioutil.ReadFile(bankPath)
	if err != nil {
		http.Error(w, "Cannot read problem bank: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var all []Problem
	if err := json.Unmarshal(raw, &all); err != nil {
		http.Error(w, "Invalid problem bank: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var tasks []Problem
	for _, p := range all {
		if p.Topic == topic {
			tasks = append(tasks, p)
		}
	}

	if len(tasks) == 0 {
		http.NotFound(w, r)
		return
	}

	data := TasksData{
		Subject: subject,
		Topic:   topic,
		Topics:  []string{},
		Tasks:   tasks,
		Checked: false,
	}

	// === Обработка POST ===
	if r.Method == http.MethodPost {
		data.Checked = true
		data.Results = make(map[int]bool)
		data.UserAns = make(map[int]string)

		for _, task := range tasks {
			key := fmt.Sprintf("answer_%d", task.ID)
			userInput := r.FormValue(key)
			data.UserAns[task.ID] = userInput

			// Проверка: есть ли совпадение с одним из допустимых ответов
			correct := false
			for _, ans := range task.Answers {
				if userInput == ans {
					correct = true
					break
				}
			}
			data.Results[task.ID] = correct
		}
	}

	title := fmt.Sprintf("Домашка: %s — %s", subject, topic)
	render(w, r, "tasks", title, data)
}

// tasksTopicsData — payload для /tasks/{subject}
type TasksTopicsData struct {
	Subject string
	Topics  []string
}

// tasksTopics показывает список тем для заданного subject
func TasksTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subject := vars["subject"] // "math" или "infa"

	// Читаем bank.json
	bankPath := filepath.Join("static", "courses", subject, "problems", "bank.json")
	raw, err := ioutil.ReadFile(bankPath)
	if err != nil {
		http.Error(w, "Cannot read problem bank: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var all []Problem
	if err := json.Unmarshal(raw, &all); err != nil {
		http.Error(w, "Invalid problem bank: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Собираем уникальные темы
	topicsMap := make(map[string]struct{})
	for _, p := range all {
		topicsMap[p.Topic] = struct{}{}
	}
	topics := make([]string, 0, len(topicsMap))
	for t := range topicsMap {
		topics = append(topics, t)
	}
	sort.Strings(topics) // опционально, чтобы было в порядке

	data := TasksTopicsData{Subject: subject, Topics: topics}
	title := fmt.Sprintf("Домашка: %s — темы", subject)

	render(w, r, "tasks_topics", title, data)
}

// tasksSubjectsData
type TasksSubjectsData struct {
	Subjects []string
}

func TasksSubjects(w http.ResponseWriter, r *http.Request) {
	data := TasksSubjectsData{Subjects: []string{"math", "infa"}}
	render(w, r, "tasks_subjects", "Домашние задания", data)
}
