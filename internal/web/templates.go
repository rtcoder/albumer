package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
)

//go:embed templates/*.gohtml assets/* locales/*.json
var webFiles embed.FS

func parseTemplates() *template.Template {
	return template.Must(template.New("albumer").Funcs(template.FuncMap{
		"eq":         func(a, b any) bool { return fmt.Sprint(a) == fmt.Sprint(b) },
		"pathEscape": url.PathEscape,
	}).ParseFS(webFiles, "templates/*.gohtml"))
}

func assetsFS() fs.FS {
	assets, err := fs.Sub(webFiles, "assets")
	if err != nil {
		panic(err)
	}
	return assets
}
