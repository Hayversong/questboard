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
	http.HandleFunc(
		"/cards/move",
		handler.MoveCardHandler,
	)
	http.HandleFunc(
		"/projects",
		handler.CreateProjectHandler,
	)
	http.HandleFunc(
		"/projects/delete",
		handler.DeleteProjectHandler,
	)
	http.HandleFunc(
		"/projects/rename",
		handler.RenameProjectHandler,
	)

	http.HandleFunc(
		"/cards/delete",
		handler.DeleteCardHandler,
	)

	http.HandleFunc(
		"/cards/update",
		handler.UpdateCardHandler,
	)

	log.Println("Servidor rodando em http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
