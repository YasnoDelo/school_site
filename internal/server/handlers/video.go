// internal/server/handlers/video.go
package handlers

import "net/http"

type VideoData struct {
	VideoURL string
}

func VideoPage(w http.ResponseWriter, r *http.Request) {
	data := VideoData{VideoURL: "/static/videos/lesson1.mp4"}
	render(w, r, "video", "Видеоурок", data)
}
