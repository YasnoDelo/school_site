package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// VideoPageData — данные, которые передаются в шаблон
type VideoPageData struct {
	VideoURL string // URL до видеофайла
}

func VideoPage(w http.ResponseWriter, r *http.Request) {
	// Формируем URL к видео (относительно корня сайта)
	data := VideoPageData{
		VideoURL: "/static/videos/lesson1.mp4",
	}

	// Рендерим шаблоны: base.html + video.html
	base := filepath.Join(TemplatesDir, "base.html")
	vid := filepath.Join(TemplatesDir, "video.html")
	tmpl := template.Must(template.ParseFiles(base, vid))

	// Выполняем рендеринг, передаём data
	tmpl.ExecuteTemplate(w, "base", data)
}
