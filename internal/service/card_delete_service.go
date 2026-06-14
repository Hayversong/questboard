package service

import (
	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func DeleteCard(projectID string, cardID string) error {

	if projectID == "" {
		return ErrProjectIDRequired
	}

	if cardID == "" {
		return ErrCardIDRequired
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		return err
	}

	for p := range projects {

		if projects[p].ID != projectID {
			continue
		}

		var updatedCards []model.Card
		cardDeleted := false

		for _, card := range projects[p].Cards {

			if card.ID == cardID {
				AddActivity(&projects[p], "🗑️ Quest removida: "+card.Title)
				cardDeleted = true
				continue
			}

			updatedCards = append(updatedCards, card)
		}

		if !cardDeleted {
			return ErrCardNotFound
		}

		projects[p].Cards = updatedCards

		return storage.SaveProjects(projects)
	}

	return ErrProjectNotFound
}
