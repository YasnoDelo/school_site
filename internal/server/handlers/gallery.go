// internal/server/handlers/gallery.go
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type GalleryData struct {
	Images []string
}

func Gallery(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		if _, err := os.Stat(filepath.Join(dir, "static", "images", "gallery")); err == nil {
			break
		}
		dir = filepath.Dir(dir)
	}

	matches, err := filepath.Glob(filepath.Join(dir, "static", "images", "gallery", "*"))
	if err != nil {
		http.Error(w, "Cannot read gallery: "+err.Error(), 500)
		return
	}
	var imgs []string
	for _, p := range matches {
		rel := strings.TrimPrefix(p, dir)
		imgs = append(imgs, filepath.ToSlash(rel))
	}
	render(w, r, "gallery", "Галерея", GalleryData{Images: imgs})
}
