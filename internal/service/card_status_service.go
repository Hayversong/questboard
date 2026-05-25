package service

import (
	"errors"

	"github.com/Hayversong/questboard/internal/storage"
)

func UpdateCardStatus(
	projectID string,
	cardID string,
	status string,
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

		for c := range projects[p].Cards {

			card :=
				&projects[p].Cards[c]

			if card.ID != cardID {
				continue
			}

			card.Status =
				status

			return storage.SaveProjects(
				projects,
			)
		}
	}

	return errors.New(
		"card não encontrado",
	)
}
