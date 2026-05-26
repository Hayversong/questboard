package handler

import (
	"html/template"
	"net/http"
	"sort"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/service"
	"github.com/Hayversong/questboard/internal/storage"
)

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

func ProjectDetailHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.URL.Query().Get(
		"id",
	)

	if id == "" {

		http.Error(
			w,
			"ID obrigatório",
			http.StatusBadRequest,
		)

		return
	}

	projects, err := storage.LoadProjects()

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	for _, project := range projects {

		if project.ID == id {

			sort.Slice(
				project.Cards,
				func(i, j int) bool {

					return project.Cards[i].Order <
						project.Cards[j].Order
				},
			)

			tmpl, err := template.ParseFiles(
				"web/templates/project.html",
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
				project,
			)

			return
		}
	}

	http.Error(
		w,
		"Projeto não encontrado",
		http.StatusNotFound,
	)
}
