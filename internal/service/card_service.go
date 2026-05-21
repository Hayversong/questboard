package service

import (
	"fmt"
	"time"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func CreateCard(
	projectID string,
	title string,
	description string,
	rarity string,
) error {

	projects, err := storage.LoadProjects()

	if err != nil {
		return err
	}

	for i := range projects {

		if projects[i].ID != projectID {
			continue
		}

		card := model.Card{
			ID: fmt.Sprintf(
				"%d",
				time.Now().UnixNano(),
			),

			Title: title,

			Description: description,

			Status: "backlog",

			Rarity: rarity,
		}

		projects[i].Cards = append(
			projects[i].Cards,
			card,
		)

		break
	}

	return storage.SaveProjects(
		projects,
	)
}

func MoveCard(
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
		}
	}

	return storage.SaveProjects(
		projects,
	)
}
