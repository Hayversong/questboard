package handler

import (
	"net/http"

	"github.com/Hayversong/questboard/internal/service"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")

	err := service.CreateProject(name)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
