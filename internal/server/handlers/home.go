// internal/server/handlers/home.go
package handlers

import "net/http"

func Home(w http.ResponseWriter, r *http.Request) {
	render(w, r, "home", "Главная", nil)
}
