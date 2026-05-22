package handler

import (
	"html/template"
	"net/http"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/service"
	"github.com/Hayversong/questboard/internal/storage"
)

// HomeHandler renderiza a lista de projetos
func HomeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	projects, err := storage.LoadProjects()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	stats, err := service.DashboardStats()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	data := struct {
		Projects []model.Project
		Stats    model.Dashboard
	}{
		Projects: projects,
		Stats:    stats,
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

	tmpl.Execute(
		w,
		data,
	)
}
