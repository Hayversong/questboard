package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func CreateProject(name string) error {

	name = strings.TrimSpace(name)

	if name == "" {
		return ErrProjectNameRequired
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		return err
	}

	project := model.Project{
		ID: fmt.Sprintf(
			"%d",
			time.Now().UnixNano(),
		),
		Name:  name,
		Cards: []model.Card{},
	}

	projects = append(projects, project)

	return storage.SaveProjects(projects)
}

func DashboardStats() (model.Dashboard, error) {

	projects, err := storage.LoadProjects()
	if err != nil {
		return model.Dashboard{}, err
	}

	stats := model.Dashboard{}

	stats.Projects = len(projects)

	for _, project := range projects {

		stats.XP += project.XP()

		stats.Cards += len(project.Cards)

		for _, card := range project.Cards {
			if card.Status == "done" {
				stats.Done++
			}
		}
	}

	return stats, nil
}

func GetProjectByID(projectID string) (model.Project, error) {

	if projectID == "" {
		return model.Project{}, ErrProjectIDRequired
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		return model.Project{}, err
	}

	for _, project := range projects {

		if project.ID == projectID {

			sort.Slice(project.Cards, func(i, j int) bool {
				return project.Cards[i].Order < project.Cards[j].Order
			})

			return project, nil
		}
	}

	return model.Project{}, ErrProjectNotFound
}
