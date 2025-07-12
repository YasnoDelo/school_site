// internal/server/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/YasnoDelo/school_site/internal/server/handlers"
)

// AuthRequired оборачивает HTTP‑хендлер, требуя авторизации.
// Если в сессии нет ключа "user_id", делает редирект на /login.
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Берём сессию из глобального хранилища
		sess, err := handlers.SessionStore.Get(r, "session-name")
		if err != nil {
			// Не смогли получить сессию — считаем, что пользователь не в системе
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Проверяем, есть ли user_id
		userID, ok := sess.Values["user_id"]
		if !ok || userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Всё ок — передаём дальше
		next.ServeHTTP(w, r)
	})
}
