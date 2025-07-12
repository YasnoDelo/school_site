// internal/server/handlers/auth.go
package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Register
type RegisterForm struct {
	Error string
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	render(w, r, "register", "Регистрация", RegisterForm{})
}

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := DB.Exec(
		"INSERT INTO users (username, password) VALUES ($1,$2)",
		username, string(hash),
	)
	if err != nil {
		render(w, r, "register", "Регистрация", RegisterForm{Error: "User exists or DB error"})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login
type LoginForm struct {
	Error string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	render(w, r, "login", "Вход", LoginForm{})
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var id int
	var hash string
	err := DB.QueryRow(
		"SELECT id, password FROM users WHERE username=$1",
		username,
	).Scan(&id, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		render(w, r, "login", "Вход", LoginForm{Error: "Invalid credentials"})
		return
	}

	sess, _ := SessionStore.Get(r, "session-name")
	sess.Values["user_id"] = id
	sess.Save(r, w)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionStore.Get(r, "session-name")
	delete(sess.Values, "user_id")
	// Просим браузер удалять cookie сразу
	sess.Options.MaxAge = -1
	// Записываем заголовок Set-Cookie
	if err := sess.Save(r, w); err != nil {
		http.Error(w, "Logout error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Profile
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	render(w, r, "profile", "Личный кабинет", nil)
}
