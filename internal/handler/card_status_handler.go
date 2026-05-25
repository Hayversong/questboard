package handler

import (
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func UpdateCardStatusHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {

		http.Error(
			w,
			"método inválido",
			http.StatusMethodNotAllowed,
		)

		return
	}

	projectID :=
		r.FormValue(
			"project_id",
		)

	cardID :=
		r.FormValue(
			"card_id",
		)

	status :=
		r.FormValue(
			"status",
		)

	err :=
		service.UpdateCardStatus(
			projectID,
			cardID,
			status,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			500,
		)

		return
	}

	w.WriteHeader(
		http.StatusOK,
	)
}
