// internal/server/handlers/render.go
package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// render теперь принимает ещё title string
func render(w http.ResponseWriter, r *http.Request, page, title string, payload interface{}) {
	sess, _ := SessionStore.Get(r, "session-name")
	uid, ok := sess.Values["user_id"].(int)

	vd := ViewData{Title: title}
	if ok {
		var u User
		if err := DB.QueryRow(
			"SELECT id, username FROM users WHERE id=$1", uid,
		).Scan(&u.ID, &u.Username); err == nil {
			vd.IsAuth = true
			vd.Username = u.Username
		}
	}

	vd.Data = payload

	base := filepath.Join(TemplatesDir, "base.html")
	tpl := filepath.Join(TemplatesDir, page+".html")
	tmpl := template.Must(template.ParseFiles(base, tpl))
	if err := tmpl.ExecuteTemplate(w, "base", vd); err != nil {
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
	}
}
