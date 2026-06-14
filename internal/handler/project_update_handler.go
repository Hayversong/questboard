package handler

import (
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func RenameProjectHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.FormValue("project_id")
	name := r.FormValue("name")

	err := service.RenameProject(projectID, name)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
