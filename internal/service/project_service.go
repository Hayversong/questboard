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
