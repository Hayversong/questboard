package handler

import (
	"fmt"
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func CreateCardHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Método não permitido",
			http.StatusMethodNotAllowed,
		)

		return
	}

	err := r.ParseForm()

	if err != nil {

		http.Error(
			w,
			"Erro ao processar formulário",
			http.StatusBadRequest,
		)

		return
	}

	projectID := r.FormValue(
		"project_id",
	)

	title := r.FormValue(
		"title",
	)

	description := r.FormValue(
		"description",
	)

	rarity := r.FormValue(
		"rarity",
	)

	if projectID == "" || title == "" {

		http.Error(
			w,
			"Projeto e título obrigatórios",
			http.StatusBadRequest,
		)

		return
	}

	err = service.CreateCard(
		projectID,
		title,
		description,
		rarity,
	)

	if err != nil {

		http.Error(
			w,
			"Erro ao criar card",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/project?id=%s",
			projectID,
		),
		http.StatusSeeOther,
	)
}

func MoveCardHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"Método não permitido",
			http.StatusMethodNotAllowed,
		)

		return
	}

	projectID := r.FormValue(
		"project_id",
	)

	cardID := r.FormValue(
		"card_id",
	)

	err := service.MoveCard(
		projectID,
		cardID,
	)

	if err != nil {

		http.Error(
			w,
			"Erro ao mover card",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"/project?id=%s",
			projectID,
		),
		http.StatusSeeOther,
	)
}
