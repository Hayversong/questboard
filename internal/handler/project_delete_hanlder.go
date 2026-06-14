package handler

import (
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.FormValue("project_id")

	err := service.DeleteProject(projectID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
