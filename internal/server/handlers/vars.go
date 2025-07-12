package handlers

import (
	"database/sql"

	"github.com/gorilla/sessions"
)

// эти переменные устанавливаются из server.NewServerMux
var (
	TemplatesDir string
	DataDir      string
)

var (
	DB           *sql.DB
	SessionStore *sessions.CookieStore
)
