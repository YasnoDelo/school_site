package main

import (
	"log"
	"net/http"

	"github.com/YasnoDelo/school_site/internal/server"
	"github.com/YasnoDelo/school_site/internal/server/config"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	router := mux.NewRouter()
	muxServer := server.NewServerMux(cfg, router)

	addr := ":" + cfg.Port
	log.Printf("Starting HTTP server on %s (env=%s)", addr, cfg.Env)
	log.Fatal(http.ListenAndServe(addr, muxServer))
}
