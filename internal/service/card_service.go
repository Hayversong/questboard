package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func CreateCard(
	projectID string,
	title string,
	description string,
	rarity string,
	deadline string,
) error {

	if projectID == "" {
		return ErrProjectIDRequired
	}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		return ErrCardTitleRequired
	}

	if rarity == "" {
		rarity = "common"
	}

	if !isValidRarity(rarity) {
		return ErrInvalidRarity
	}

	if !isValidDeadline(deadline) {
		return ErrInvalidDeadline
	}

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
			Title:       title,
			Description: description,
			Status:      "backlog",
			Rarity:      rarity,
			Deadline:    deadline,
			Order:       len(projects[i].Cards),
		}

		projects[i].Cards = append(projects[i].Cards, card)

		AddActivity(&projects[i], "⭐ Quest criada: "+title)

		return storage.SaveProjects(projects)
	}

	return ErrProjectNotFound
}

func MoveCard(projectID string, cardID string) error {

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

		for c := range projects[p].Cards {

			card := &projects[p].Cards[c]

			if card.ID != cardID {
				continue
			}

			switch card.Status {

			case "backlog":
				card.Status = "doing"
				AddActivity(&projects[p], "🚀 Quest iniciada: "+card.Title)

			case "doing":
				card.Status = "done"
				AddActivity(&projects[p], "✅ Quest concluída: "+card.Title)

			default:
				return nil
			}

			return storage.SaveProjects(projects)
		}

		return ErrCardNotFound
	}

	return ErrProjectNotFound
}
