package main

import (
	"fmt"
	"net/http"

	"github.com/AndreiPatriota/go-sse/internal/routes"
)

func main() {
	r := routes.InitRouter()

	port := 3853
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	fmt.Printf("Server is running on port %d\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}