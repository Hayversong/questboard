package handler

import (
	"fmt"
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func DeleteCardHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.FormValue("project_id")
	cardID := r.FormValue("card_id")

	err := service.DeleteCard(projectID, cardID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/project?id=%s", projectID), http.StatusSeeOther)
}
