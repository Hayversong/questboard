package service

import (
	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func DeleteProject(
	projectID string,
) error {

	projects, err := storage.LoadProjects()

	if err != nil {
		return err
	}

	var updated []model.Project

	for _, p := range projects {

		if p.ID != projectID {
			updated = append(
				updated,
				p,
			)
		}
	}

	return storage.SaveProjects(
		updated,
	)
}
