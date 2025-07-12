// internal/server/handlers/vars.go
package handlers

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

var (
	DB           *sql.DB
	SessionStore *sessions.CookieStore
	TemplatesDir string
)

type User struct {
	ID       int
	Username string
}

// ViewData — единый контейнер для всех шаблонов
type ViewData struct {
	Title    string // заголовок страницы
	IsAuth   bool
	Username string
	Data     interface{}
}
