// internal/server/handlers/materials.go
package handlers

import "net/http"

type MaterialsData struct {
	Files []string
}

func Materials(w http.ResponseWriter, r *http.Request) {
	render(w, r, "materials", "Материалы", nil)
}
