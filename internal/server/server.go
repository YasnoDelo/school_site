package server

import (
	"net/http"
	"os"
	"path/filepath"

	"log"

	"github.com/YasnoDelo/school_site/internal/server/config"
	"github.com/YasnoDelo/school_site/internal/server/handlers"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func NewServerMux(cfg *config.Config, router *mux.Router) http.Handler {
	cwd, err := os.Getwd()
	if err != nil {
		panic("cannot get working dir: " + err.Error())
	}
	// cwd == .../Best_HTTP_server_ever/cmd/server
	projectRoot := filepath.Dir(filepath.Dir(cwd))
	templatesDir := filepath.Join(projectRoot, "internal", "templates")

	// Устанавливаем глобальную переменную в пакете handlers
	handlers.TemplatesDir = templatesDir

	// Middleware, статика, роуты как было
	router.Use(loggingMiddleware)
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.
		PathPrefix("/img/").
		Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	router.HandleFunc("/", handlers.Home).Methods("GET")
	router.HandleFunc("/courses", handlers.Courses).Methods("GET")
	router.HandleFunc("/materials", handlers.Materials).Methods("GET")

	return router
}
