package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// TeamMember — данные для одного человека
type TeamMember struct {
	Name  string
	Role  string
	Bio   string
	Photo string // относительный путь от /static, например "images/team/alice.jpg"
}

// AboutData — payload для шаблона
type AboutData struct {
	Members []TeamMember
}

// About — хендлер страницы /about
func About(w http.ResponseWriter, r *http.Request) {
	// Папка с фотками относительно рабочей директории
	teamDir := filepath.Join("static", "images", "team")

	// Попробуем прочитать файлы в каталоге и автоматом сформировать список
	files, err := os.ReadDir(teamDir)
	var members []TeamMember
	if err == nil {
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
				continue
			}
			// Берём имя файла без расширения как имя человека (например alice -> Alice)
			base := strings.TrimSuffix(name, ext)
			displayName := strings.Title(strings.ReplaceAll(base, "_", " "))
			members = append(members, TeamMember{
				Name:  displayName,
				Role:  "", // можно заполнить вручную ниже
				Bio:   "",
				Photo: filepath.ToSlash(filepath.Join("images", "team", name)), // относительный путь внутри /static
			})
		}
	}

	// Если автоматический список пуст — можно указать людей вручную
	if len(members) == 0 {
		members = []TeamMember{
			{Name: "Никита", Role: "Преподаватель по информатике", Bio: "Что-то там", Photo: "images/team/Mikito.jpg"},
			{Name: "Григорий", Role: "Преподаватель по математике", Bio: "Что-то там", Photo: "images/team/Grigory.jpg"},
		}
	}

	data := AboutData{Members: members}
	render(w, r, "about", "О нас", data)
}
