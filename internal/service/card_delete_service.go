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

	projects, err :=
		storage.LoadProjects()

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

				AddActivity(
					&projects[p],
					"🗑️ Quest removida: "+
						card.Title,
				)

				cardDeleted = true

				continue
			}

			updatedCards =
				append(
					updatedCards,
					card,
				)
		}

		if !cardDeleted {

			return errors.New(
				"card não encontrado",
			)
		}

		projects[p].Cards =
			updatedCards

		return storage.SaveProjects(
			projects,
		)
	}

	return errors.New(
		"projeto não encontrado",
	)
}
