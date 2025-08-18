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

	members = []TeamMember{
		{Name: "Никита", Role: "Преподаватель по информатике", Bio: "В МФТИ (где я учусь уже 4 года) есть негласное правило: «Разобрался сам - объясни другим». В том числе из-за этого дивиза я уже 4 года готовлю к государственным экзаменам и делаю скучное интересным, а непонятное - очевидным. \n\nНа уроках я решаю 2 задачи. Во-первых мотивирую ученика, чтобы он сам хотел разобраться. Во-вторых рассказываю материал так, чтобы подопечный сам мог (хотел) объяснить услышанное кому угодно", Photo: "images/team/Mikito.jpg"},
		{Name: "Григорий", Role: "Преподаватель по математике", Bio: "Уже 4 года я преподаю математику и физику, помогая школьникам готовиться к ОГЭ и ЕГЭ. Учусь в МГСУ на 4-м курсе с отличием, ранее окончил ЗФТШ при МФТИ и становился победителем олимпиад 2-го урвня. \n\nВ обучении делаю упор на то, чтобы ребёнок сам находил решения, а я лишь направляю и помогаю шаг за шагом. Такой подход не только повышает результаты, но и формирует уверенность в собственных силах.", Photo: "images/team/Grigory.jpg"},
	}

	data := AboutData{Members: members}
	render(w, r, "about", "О нас", data)
}
