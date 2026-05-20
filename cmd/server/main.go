package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Hayversong/questboard/internal/storage"
)

func home(w http.ResponseWriter, r *http.Request) {

	projects, err := storage.LoadProjects()

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	tmpl, err := template.ParseFiles(
		"web/templates/home.html",
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = tmpl.Execute(
		w,
		projects,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

func main() {

	http.HandleFunc("/", home)

	log.Println("Servidor rodando em :8080")

	err := http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}
}
