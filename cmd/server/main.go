package main

import (
	"log"
	"net/http"

	"github.com/Hayversong/questboard/internal/handler"
)

func main() {
	// Rotas de visualização
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/project", handler.ProjectDetailHandler)

	// Rotas de ação
	http.HandleFunc("/cards", handler.CreateCardHandler)

	log.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
