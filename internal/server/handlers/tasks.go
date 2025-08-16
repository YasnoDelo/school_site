package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

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

// normalize — убирает пробелы, обрезает $, приводит к нижнему регистру
func normalize(s string) string {
	s = strings.TrimSpace(s)
	// убрать окружающие $ (формулы) и фигурные скобки
	s = strings.Trim(s, "$ \t\n\r")
	s = strings.Trim(s, "{}")
	// заменить NBSP на обычный пробел
	s = strings.ReplaceAll(s, "\u00A0", " ")
	// collapse внутренние пробелы
	s = strings.Join(strings.Fields(s), " ")
	return strings.ToLower(s)
}

// tryParseNumber — пытается получить float64 из строки.
// Поддерживает "1.23", "1,23" и простые дроби "3/4".
// Возвращает (value, true) при успехе.
func tryParseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}

	// поддержка дробей вида a/b (включая знаки и запятые)
	if strings.Contains(s, "/") {
		// убираем пробелы вокруг '/'
		parts := strings.SplitN(s, "/", 2)
		if len(parts) == 2 {
			a := strings.ReplaceAll(strings.TrimSpace(parts[0]), ",", ".")
			b := strings.ReplaceAll(strings.TrimSpace(parts[1]), ",", ".")
			// попытаемся через big.Rat
			ra := new(big.Rat)
			rb := new(big.Rat)
			okA := false
			okB := false
			if _, ok := ra.SetString(a); ok {
				okA = true
			} else {
				// как запас: если a - целое/float
				if fa, err := strconv.ParseFloat(a, 64); err == nil {
					ra.SetFloat64(fa)
					okA = true
				}
			}
			if _, ok := rb.SetString(b); ok {
				okB = true
			} else {
				if fb, err := strconv.ParseFloat(b, 64); err == nil {
					rb.SetFloat64(fb)
					okB = true
				}
			}
			if okA && okB {
				r := new(big.Rat).Quo(ra, rb)
				f, _ := r.Float64()
				return f, true
			}
		}
	}

	// заменяем запятую на точку и пробуем float
	s2 := strings.ReplaceAll(s, ",", ".")
	if f, err := strconv.ParseFloat(s2, 64); err == nil {
		return f, true
	}
	return 0, false
}

// answersEqual — true если ответы совпадают по строке (без учёта регистра/пробелов)
// или являются численно равными (учитывая , или . и дроби).
func answersEqual(userAns, correctAns string) bool {
	nu := normalize(userAns)
	nc := normalize(correctAns)

	// 1) точное строковое совпадение (case-insensitive)
	if nu == nc {
		return true
	}

	// 2) попробуем числовое сравнение
	fu, okU := tryParseNumber(nu)
	fc, okC := tryParseNumber(nc)
	if okU && okC {
		const eps = 1e-9
		if math.Abs(fu-fc) <= eps {
			return true
		}
	}

	// 3) ещё попытка: убрать пробелы и сравнить (полезно для latex без пробелов)
	if strings.ReplaceAll(nu, " ", "") == strings.ReplaceAll(nc, " ", "") {
		return true
	}

	return false
}
