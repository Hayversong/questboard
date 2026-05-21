package handler

import (
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func CreateProjectHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Método inválido",
			http.StatusMethodNotAllowed,
		)

		return
	}

	name := r.FormValue(
		"name",
	)

	if name == "" {

		http.Error(
			w,
			"Nome obrigatório",
			http.StatusBadRequest,
		)

		return
	}

	err := service.CreateProject(
		name,
	)

	if err != nil {

		http.Error(
			w,
			"Erro ao criar",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		"/",
		http.StatusSeeOther,
	)
}
