package main

import (
	"log"
	"net/http"

	"github.com/YasnoDelo/school_site/internal/server"
	"github.com/YasnoDelo/school_site/internal/server/config"
)

func main() {
	// 1. Загружаем конфиг (.env или OS‑env)
	cfg := config.Load()
	log.Println("DatabaseDSN:", cfg.DatabaseDSN)

	// 2. Создаём весь сервер (DB, миграции, сессии, маршруты и т.д.)
	handler := server.NewServerMux(cfg)

	// 3. Запускаем HTTP‑сервер
	addr := ":" + cfg.Port
	log.Printf("Starting server on %s (env=%s)", addr, cfg.Env)
	log.Fatal(http.ListenAndServe(addr, handler))
}
