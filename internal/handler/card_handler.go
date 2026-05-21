package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	projectID := r.FormValue("project_id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	if projectID == "" || title == "" {
		http.Error(w, "ID do projeto e título são obrigatórios", http.StatusBadRequest)
		return
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		http.Error(w, "Erro ao carregar projetos", http.StatusInternalServerError)
		return
	}

	// Encontrar o projeto e adicionar o card
	found := false
	for i := range projects {
		if projects[i].ID == projectID {
			newCard := model.Card{
				ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
				Title:       title,
				Description: description,
				Status:      "backlog",
			}
			projects[i].Cards = append(projects[i].Cards, newCard)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Projeto não encontrado", http.StatusNotFound)
		return
	}

	err = storage.SaveProjects(projects)
	if err != nil {
		http.Error(w, "Erro ao salvar projeto", http.StatusInternalServerError)
		return
	}

	// Redireciona de volta para a página do projeto
	http.Redirect(w, r, fmt.Sprintf("/project?id=%s", projectID), http.StatusSeeOther)
}
