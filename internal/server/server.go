package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/YasnoDelo/school_site/internal/server/config"
	"github.com/YasnoDelo/school_site/internal/server/handlers"
	"github.com/gorilla/mux"
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

// findProjectRoot поднимается от cwd вверх, пока не найдёт папку .git или просто на N уровней
func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get cwd: %v", err)
	}
	// пробуем найти папку data рядом с internal, либо поднимаемся N раз
	dir := cwd
	for i := 0; i < 4; i++ {
		cand := filepath.Join(dir, "data")
		if info, err := os.Stat(cand); err == nil && info.IsDir() {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	log.Fatalf("could not locate project root from %s", cwd)
	return ""
}

func NewServerMux(cfg *config.Config, router *mux.Router) http.Handler {
	// Ищем папку с шаблонами автоматически
	templatesDir := findTemplatesDir()
	handlers.TemplatesDir = templatesDir

	// Middleware
	router.Use(loggingMiddleware)

	// Статика
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.
		PathPrefix("/img/").
		Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Роуты
	router.HandleFunc("/", handlers.Home).Methods("GET")
	router.HandleFunc("/courses", handlers.Courses).Methods("GET")
	router.HandleFunc("/materials", handlers.Materials).Methods("GET")
	router.HandleFunc("/homework", handlers.Homework).Methods("GET", "POST")
	router.HandleFunc("/gallery", handlers.Gallery).Methods("GET")
	router.HandleFunc("/video", handlers.VideoPage).Methods("GET")

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[REQUEST] %s %s\n", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	return router
}

// loggingMiddleware — простой логер
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
