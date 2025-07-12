package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

// User — модель для примера
type User struct {
	ID       int
	Username string
	Password string // хранится bcrypt-хеш
}

// renderTemplate — удобная функция
func render(w http.ResponseWriter, name string, data interface{}) {
	base := filepath.Join(TemplatesDir, "base.html")
	tpl := filepath.Join(TemplatesDir, name+".html")
	tmpl := template.Must(template.ParseFiles(base, tpl))
	tmpl.ExecuteTemplate(w, "base", data)
}

// SignUp — регистрация
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "signup", nil)
		return
	}
	// POST
	username := r.FormValue("username")
	password := r.FormValue("password")
	// хешируем
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}
	// сохраняем в БД
	_, err = DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hash))
	if err != nil {
		http.Error(w, "DB error", 500)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login — вход
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "login", nil)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	// ищем в БД
	var user User
	err := DB.QueryRow("SELECT id, password FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}
	// сравниваем хеш
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}
	// создаём сессию
	sess, _ := SessionStore.Get(r, "session-name")
	sess.Values["user_id"] = user.ID
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout — выход
func Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionStore.Get(r, "session-name")
	delete(sess.Values, "user_id")
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
