package service

import (
	"fmt"
	"time"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func CreateProject(
	name string,
) error {

	projects, err := storage.LoadProjects()

	if err != nil {
		return err
	}

	project := model.Project{
		ID: fmt.Sprintf(
			"%d",
			time.Now().UnixNano(),
		),

		Name: name,

		Cards: []model.Card{},
	}

	projects = append(
		projects,
		project,
	)

	return storage.SaveProjects(
		projects,
	)
}

func DashboardStats() (model.Dashboard, error) {

	projects, err := storage.LoadProjects()

	if err != nil {
		return model.Dashboard{}, err
	}

	stats := model.Dashboard{}

	stats.Projects = len(
		projects,
	)

	for _, project := range projects {

		stats.XP += project.XP()

		stats.Cards += len(
			project.Cards,
		)

		for _, card := range project.Cards {

			if card.Status == "done" {
				stats.Done++
			}
		}
	}

	return stats, nil
}