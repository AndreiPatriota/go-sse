package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/home")
	w.WriteHeader(http.StatusFound)
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("web/views/_layout.html" ,"web/views/home.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "home", nil)
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("web/views/_layout.html" ,"web/views/app.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "app", nil)
}

func GetSseStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	notificaFim := r.Context().Done()

	for {
		select {
		case <-notificaFim:
			log.Println("Client closed connection")
			return
		case t := <-tick.C:
			msg := mensagemSSE{
				Ok:       true,
				Mensagem: fmt.Sprintf("Current time is %s", t.Format(time.RFC3339)),
			}
			jsonPayload, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", jsonPayload)
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
			flusher.Flush()
		}
	}
}