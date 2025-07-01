package main

import (
	"log"
	"net/http"

	"github.com/YasnoDelo/school_site/internal/server"
	"github.com/YasnoDelo/school_site/internal/server/config"
	"github.com/gorilla/mux"
)

func main() {
	// 1. Загрузка конфигурации (.env или OS‑env)
	cfg := config.Load()

	// 2. Создаем Gorilla Mux
	router := mux.NewRouter()

	// 3. Передаем router и cfg в NewServerMux,
	//    который внутри настроит статику, middleware и хендлеры
	handler := server.NewServerMux(cfg, router)

	// 4. Запускаем HTTP‑сервер
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s (env=%s)", addr, cfg.Env)
	log.Fatal(http.ListenAndServe(addr, handler))
}
