package handler

import (
	"html/template"
	"net/http"

	"github.com/Hayversong/questboard/internal/storage"
)

// HomeHandler renderiza a lista de projetos
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	projects, err := storage.LoadProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, projects)
}

// ProjectDetailHandler renderiza os detalhes de um projeto específico (Kanban)
func ProjectDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID do projeto é obrigatório", http.StatusBadRequest)
		return
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var projectFound interface{}
	for _, p := range projects {
		if p.ID == id {
			projectFound = p
			break
		}
	}

	if projectFound == nil {
		http.Error(w, "Projeto não encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/project.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, projectFound)
}
