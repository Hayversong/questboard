package main

import (
	"log"
	"net/http"

	"github.com/Hayversong/questboard/internal/handler"
)

func main() {

	// Arquivos estáticos
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(
				http.Dir(
					"web/static",
				),
			),
		),
	)

	// Visualização
	http.HandleFunc(
		"/",
		handler.HomeHandler,
	)

	http.HandleFunc(
		"/project",
		handler.ProjectDetailHandler,
	)

	// Cards
	http.HandleFunc(
		"/cards",
		handler.CreateCardHandler,
	)

	http.HandleFunc(
		"/cards/move",
		handler.MoveCardHandler,
	)

	http.HandleFunc(
		"/cards/delete",
		handler.DeleteCardHandler,
	)

	http.HandleFunc(
		"/cards/update",
		handler.UpdateCardHandler,
	)

	// Projetos
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
		"/cards/status",
		handler.UpdateCardStatusHandler,
	)

	log.Println(
		"Servidor rodando em http://localhost:8080",
	)

	log.Fatal(
		http.ListenAndServe(
			":8080",
			nil,
		),
	)
}
