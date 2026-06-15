package service

import (
	"github.com/Hayversong/questboard/internal/storage"
)

type CardOrder struct {
	ID     string `json:"id"`
	Order  int    `json:"order"`
	Status string `json:"status"`
}

func ReorderCards(projectID string, cards []CardOrder) error {

	if projectID == "" {
		return ErrProjectIDRequired
	}

	if len(cards) == 0 {
		return ErrCardsListRequired
	}

	for _, item := range cards {
		if item.ID == "" {
			return ErrCardIDRequired
		}
		if item.Order < 0 {
			return ErrInvalidOrder
		}
		if item.Status != "" && !isValidStatus(item.Status) {
			return ErrInvalidStatus
		}
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		return err
	}

	for p := range projects {

		if projects[p].ID != projectID {
			continue
		}

		for _, item := range cards {
			for c := range projects[p].Cards {
				if projects[p].Cards[c].ID == item.ID {
					projects[p].Cards[c].Order = item.Order
					if item.Status != "" {
						projects[p].Cards[c].Status = item.Status
					}
				}
			}
		}

		return storage.SaveProjects(projects)
	}

	return ErrProjectNotFound
}
