package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hayversong/questboard/internal/handler"
)

func main() {

	// Arquivos estaticos
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

	http.HandleFunc(
		"/healthz",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Metodo invalido", http.StatusMethodNotAllowed)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		},
	)

	// Visualizacao
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

	http.HandleFunc(
		"/cards/reorder",
		handler.ReorderCardsHandler,
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Servidor rodando em http://localhost:%s", port)

	log.Fatal(
		server.ListenAndServe(),
	)
}
