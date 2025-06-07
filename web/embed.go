package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

//go:embed views/*.html
var templatesFS embed.FS

//go:embed public/**
var staticFS embed.FS

func RenderPager(w http.ResponseWriter,  templateName string, data any){
	templ, err := template.ParseFS(templatesFS, "views/_layout.html", templateName)
	if err != nil {
		log.Println("Erro ao carregar o template:", err)
		http.Error(w, "Erro ao carregar o template", http.StatusInternalServerError)
		return
	} 

	err = templ.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Erro ao renderizar o template", http.StatusInternalServerError)
		return
	}
}

func ReturnFS() fs.FS {
	staticContent, err := fs.Sub(staticFS, "public")
	if err != nil {
		log.Fatal("Erro ao criar sub-FS:", err)
	}
	return staticContent
}