package handler

import (
	"net/http"

	"github.com/Hayversong/questboard/internal/storage"
)

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

	projectID := r.FormValue("project_id")
	cardID := r.FormValue("card_id")

	projects, err := storage.LoadProjects()

	if err != nil {
		http.Error(
			w,
			"Erro ao carregar",
			http.StatusInternalServerError,
		)
		return
	}

	for p := range projects {

		if projects[p].ID != projectID {
			continue
		}

		for c := range projects[p].Cards {

			card := &projects[p].Cards[c]

			if card.ID != cardID {
				continue
			}

			switch card.Status {

			case "backlog":
				card.Status = "doing"

			case "doing":
				card.Status = "done"
			}

			break
		}
	}

	err = storage.SaveProjects(projects)

	if err != nil {
		http.Error(
			w,
			"Erro ao salvar",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/project?id="+projectID,
		http.StatusSeeOther,
	)
}