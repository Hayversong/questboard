package handler

import (
	"errors"
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

// writeServiceError traduz erros do service em respostas HTTP com o status correto.
func writeServiceError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, service.ErrProjectIDRequired),
		errors.Is(err, service.ErrCardIDRequired),
		errors.Is(err, service.ErrProjectNameRequired),
		errors.Is(err, service.ErrCardTitleRequired),
		errors.Is(err, service.ErrInvalidRarity),
		errors.Is(err, service.ErrInvalidStatus),
		errors.Is(err, service.ErrInvalidDeadline),
		errors.Is(err, service.ErrCardsListRequired),
		errors.Is(err, service.ErrInvalidOrder):
		status = http.StatusBadRequest

	case errors.Is(err, service.ErrProjectNotFound),
		errors.Is(err, service.ErrCardNotFound):
		status = http.StatusNotFound
	}

	http.Error(w, err.Error(), status)
}
