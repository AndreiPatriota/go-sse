package routes

import (
	"net/http"

	"github.com/AndreiPatriota/go-sse/internal/handlers"
	"github.com/AndreiPatriota/go-sse/web"
	"github.com/go-chi/chi/v5"
)


func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	staticContent := web.ReturnFS()
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))

	r.Get("/", handlers.GetIndex)
	r.Get("/home", handlers.GetHome)
	r.Get("/app", handlers.GetApp)
	r.Get("/sse-stream", handlers.GetSseStream)

	return r
}