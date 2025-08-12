package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldhtml "github.com/yuin/goldmark/renderer/html"
)

// mdToHTML конвертирует Markdown в безопасный HTML (template.HTML).
// Поддерживает GFM (включая таблицы).
func mdToHTML(markdown string) template.HTML {
	if markdown == "" {
		return template.HTML("")
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
		),
		goldmark.WithRendererOptions(
			goldhtml.WithUnsafe(), // разрешаем "raw" HTML внутри md (будем потом очищать)
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		log.Printf("md convert error: %v", err)
		return template.HTML("") // пусто при ошибке
	}
	rawHTML := buf.String()

	// Очистка через bluemonday — UGCPolicy + разрешить таблицы и некоторые атрибуты
	p := bluemonday.UGCPolicy()
	p.AllowElements("table", "thead", "tbody", "tfoot", "tr", "th", "td", "colgroup", "col")
	p.AllowAttrs("src", "alt", "title", "width", "height").OnElements("img")
	p.AllowAttrs("class", "style").OnElements("table", "th", "td")
	p.AllowAttrs("href", "title", "target").OnElements("a")

	safe := p.Sanitize(rawHTML)
	return template.HTML(safe)
}
