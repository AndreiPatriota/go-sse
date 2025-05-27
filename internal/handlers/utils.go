package handlers


type mensagemSSE struct {
	Ok bool `json:"ok"`
	Mensagem string `json:"mensagem"`
}