package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"github.com/YasnoDelo/school_site/internal/server/config"
	"github.com/YasnoDelo/school_site/internal/server/handlers"
	"github.com/YasnoDelo/school_site/internal/server/middleware"
)

// findTemplatesDir пытается найти папку internal/templates,
// начиная с cwd и поднимаясь вверх до трёх раз.
func findTemplatesDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working dir: %v", err)
	}

	try := cwd
	for i := 0; i < 3; i++ {
		cand := filepath.Join(try, "internal", "templates")
		if info, err := os.Stat(cand); err == nil && info.IsDir() {
			return cand
		}
		try = filepath.Dir(try)
	}

	log.Fatalf("could not locate internal/templates from cwd %s", cwd)
	return ""
}

// NewServerMux настраивает всё: БД, миграции, сессии, роуты и middleware.
func NewServerMux(cfg *config.Config) http.Handler {
	templatesDir := findTemplatesDir()
	handlers.TemplatesDir = templatesDir

	// 1) Подключаемся к Postgres
	db, err := sql.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("DB open error: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	// 2) Прогоняем миграции
	RunMigrations(db)

	// 3) Сохраняем handle на БД в пакет handlers
	handlers.DB = db

	// 4) Настраиваем Cookie‑сессии
	store := sessions.NewCookieStore([]byte(cfg.SessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   cfg.Env == "production",
	}
	handlers.SessionStore = store

	// 5) Создаём router
	r := mux.NewRouter()

	// 6) Общие middleware
	r.Use(loggingMiddleware)

	// 7) Отдача статических файлов
	//    Папка static/ лежит в корне проекта
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/img/").
		Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	// 8) Публичные маршруты
	r.HandleFunc("/", handlers.Home).Methods("GET")
	r.HandleFunc("/materials", handlers.Materials).Methods("GET")
	r.HandleFunc("/homework", handlers.Homework).Methods("GET", "POST")
	r.HandleFunc("/gallery", handlers.Gallery).Methods("GET")
	r.HandleFunc("/video", handlers.VideoPage).Methods("GET")
	r.HandleFunc("/signup", handlers.SignUp).Methods("GET", "POST")
	r.HandleFunc("/login", handlers.Login).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("POST")

	// 9) Защищённые маршруты
	r.Handle("/courses",
		middleware.AuthRequired(http.HandlerFunc(handlers.Courses)),
	).Methods("GET")

	return r
}

// loggingMiddleware — простой лог запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
