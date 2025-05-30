package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)



const (
	endpointDados = "https://api.open-meteo.com/v1/forecast"
	latItapetim = -7.3778
	logItapetim = -37.19
)

type mensagemSSE struct {
	Ok       bool   `json:"ok"`
	Mensagem string `json:"mensagem"`
	Temperatura float64 `json:"temperatura"`
	VelVento float64 `json:"velocidade_vento"`
	UmidadeRelativa int `json:"umidade_relativa"`
	Hora string `json:"hora"`
}

type weatherResponse struct {
	Latitude             float64         `json:"latitude"`
	Longitude            float64         `json:"longitude"`
	GenerationTimeMs     float64         `json:"generationtime_ms"`
	UTCOffsetSeconds     int             `json:"utc_offset_seconds"`
	Timezone             string          `json:"timezone"`
	TimezoneAbbreviation string          `json:"timezone_abbreviation"`
	Elevation            float64         `json:"elevation"`
	CurrentUnits         currentUnits    `json:"current_units"`
	Current              current         `json:"current"`
	HourlyUnits          hourlyUnits     `json:"hourly_units"`
	Hourly               hourly          `json:"hourly"`
}

type currentUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature2m string `json:"temperature_2m"`
	WindSpeed10m  string `json:"wind_speed_10m"`
}

type current struct {
	Time          string  `json:"time"` // could be time.Time with custom unmarshal
	Interval      int     `json:"interval"`
	Temperature2m float64 `json:"temperature_2m"`
	WindSpeed10m  float64 `json:"wind_speed_10m"`
}

type hourlyUnits struct {
	Time               string `json:"time"`
	Temperature2m      string `json:"temperature_2m"`
	RelativeHumidity2m string `json:"relative_humidity_2m"`
	WindSpeed10m       string `json:"wind_speed_10m"`
}

type hourly struct {
	Time               []string  `json:"time"` 
	Temperature2m      []float64 `json:"temperature_2m"`
	RelativeHumidity2m []int     `json:"relative_humidity_2m"`
	WindSpeed10m       []float64 `json:"wind_speed_10m"`
}



func getEndpointDados(lat, long float64) string {
	return endpointDados + fmt.Sprintf("?latitude=%f&longitude=%f&current=temperature_2m,wind_speed_10m&hourly=temperature_2m,relative_humidity_2m,wind_speed_10m", lat, long)
}

func getTempoLocal(lat, long float64) (*weatherResponse, error) {
	resp, err := http.Get(getEndpointDados(lat, long))
	var dados weatherResponse
	if err != nil || resp == nil {
		log.Println("erro ao fazer requisição para a API:", err)
		return nil, fmt.Errorf("erro ao fazer requisição para a API: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("erro ao obter dados:", resp.Status)
		return nil, fmt.Errorf("erro ao obter dados: %s", resp.Status)
	} 
	
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		log.Println("erro ao decodificar a resposta:", err)
		return nil, fmt.Errorf("erro ao decodificar a resposta: %v", err)
	}

	return &dados, nil
}

func sendSSEMessage(w http.ResponseWriter, msg *mensagemSSE) error {

	var event []byte = []byte("event: update\n")
	if !msg.Ok {
		event = []byte("event: error\n")
	}

	jsonPayload, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return fmt.Errorf("error marshalling JSON: %v", err)
	}

	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("Streaming não suportado!")
		return fmt.Errorf("streaming não suportado")
	} 
	
	fmt.Fprintf(w, "%sdata: %s\n\n", event, jsonPayload)
	f.Flush()
	return nil
}
