package handler

import (
	"fmt"
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func UpdateCardHandler(
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

	projectID := r.FormValue(
		"project_id",
	)

	cardID := r.FormValue(
		"card_id",
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
	deadline := r.FormValue(
		"deadline",
)

	err := service.UpdateCard(
		projectID,
		cardID,
		title,
		description,
		rarity,
		deadline,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
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
