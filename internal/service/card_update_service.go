package service

import (
	"errors"

	"github.com/Hayversong/questboard/internal/storage"
)

func UpdateCard(
	projectID string,
	cardID string,
	title string,
	description string,
	rarity string,
	deadline string,
) error {

	projects, err := storage.LoadProjects()

	if err != nil {
		return err
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

			card.Title = title
			card.Description = description
			card.Rarity = rarity
			card.Deadline = deadline

			return storage.SaveProjects(
				projects,
			)
		}
	}

	return errors.New(
		"card não encontrado",
	)
}
