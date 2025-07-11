package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/YasnoDelo/school_site/internal/server/util"
)

type GalleryData struct {
	Title  string
	Images []string
}

func Gallery(w http.ResponseWriter, r *http.Request) {
	projectRoot := util.FindProjectRoot()

	baseTmpl := filepath.Join(TemplatesDir, "base.html")
	galleryTmpl := filepath.Join(TemplatesDir, "gallery.html")
	tmpl, err := template.ParseFiles(baseTmpl, galleryTmpl)
	if err != nil {
		http.Error(w, "Error loading templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	globPath := filepath.Join(projectRoot, "static", "images", "gallery", "*")
	files, err := filepath.Glob(globPath)
	if err != nil {
		http.Error(w, "Cannot read gallery folder: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var images []string
	for _, f := range files {
		rel := strings.TrimPrefix(f, projectRoot)
		rel = filepath.ToSlash(rel)
		images = append(images, rel)
	}

	data := GalleryData{
		Title:  "Галерея",
		Images: images,
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("gallery tmpl exec error: %v", err)
		http.Error(w, "Error rendering gallery", http.StatusInternalServerError)
	}
}
