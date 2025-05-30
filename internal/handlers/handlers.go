package handlers

import (
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

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	notificaFim := r.Context().Done()

	for {
		select {
		case <-notificaFim:
			log.Println("Client closed connection")
			return
		case <-tick.C:

			var mensagem mensagemSSE
			itaDados, err := getTempoLocal(latItapetim, logItapetim)
			if err != nil || itaDados == nil {
				mensagem.Ok = false
				mensagem.Mensagem = "Erro na chamada da API"
			} else {
				mensagem.Ok = true
				mensagem.Mensagem = "Dados atualizados com sucesso"
				mensagem.Temperatura = itaDados.Current.Temperature2m
				mensagem.VelVento = itaDados.Current.WindSpeed10m
				mensagem.UmidadeRelativa = itaDados.Hourly.RelativeHumidity2m[0]
				mensagem.Hora = itaDados.Current.Time
			}
			 
			err = sendSSEMessage(w, &mensagem)
			if err != nil {
				log.Println("Error sending SSE message:", err)
				http.Error(w, "Error sending SSE message", http.StatusInternalServerError)
				return
			}
		}
	}
}