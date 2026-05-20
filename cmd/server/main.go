package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Hayversong/questboard/internal/storage"
)

func home(w http.ResponseWriter, r *http.Request) {

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

	err = tmpl.Execute(w, nil)

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

	projects, err := storage.LoadProjects()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(projects)

	log.Println("Servidor rodando em :8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
