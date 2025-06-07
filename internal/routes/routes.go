package routes

import (
	"net/http"

	"github.com/AndreiPatriota/go-sse/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/public"))))

	r.Get("/", handlers.GetIndex)
	r.Get("/home", handlers.GetHome)
	r.Get("/app", handlers.GetApp)
	r.Get("/sse-stream", handlers.GetSseStream)

	return r
}