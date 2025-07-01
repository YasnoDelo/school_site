package server

import (
	"net/http"

	"github.com/YasnoDelo/school_site/internal/server/config"
	"github.com/YasnoDelo/school_site/internal/server/handlers"
	"github.com/gorilla/mux"
)

// NewServerMux настраивает маршруты на переданном router и возвращает готовый http.Handler
func NewServerMux(cfg *config.Config, router *mux.Router) http.Handler {
	// --- Middleware (можно добавить CORS, логирование и т.д.) ---
	router.Use(loggingMiddleware)

	// --- Статические файлы ---
	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.
		PathPrefix("/img/").
		Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	// --- Динамические маршруты ---
	router.HandleFunc("/", handlers.Home).Methods("GET")
	router.HandleFunc("/courses", handlers.Courses).Methods("GET")
	router.HandleFunc("/materials", handlers.Materials).Methods("GET")

	return router
}

// Простейшее логирование запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Например, log.Printf("[%s] %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
