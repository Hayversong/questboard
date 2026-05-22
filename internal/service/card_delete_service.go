package service

import (
	"errors"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func DeleteCard(
	projectID string,
	cardID string,
) error {

	projects, err := storage.LoadProjects()

	if err != nil {
		return err
	}

	for p := range projects {

		if projects[p].ID != projectID {
			continue
		}

		var updatedCards []model.Card

		for _, card := range projects[p].Cards {

			if card.ID != cardID {
				updatedCards = append(
					updatedCards,
					card,
				)
			}
		}

		projects[p].Cards = updatedCards

		return storage.SaveProjects(
			projects,
		)
	}

	return errors.New(
		"projeto não encontrado",
	)
}
