package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

type ReorderRequest struct {
	ProjectID string              `json:"project_id"`
	Cards     []service.CardOrder `json:"cards"`
}

func ReorderCardsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	var req ReorderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = service.ReorderCards(req.ProjectID, req.Cards)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
