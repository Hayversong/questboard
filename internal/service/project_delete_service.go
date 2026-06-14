package service

import (
	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func DeleteProject(projectID string) error {

	if projectID == "" {
		return ErrProjectIDRequired
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		return err
	}

	updated := []model.Project{}
	found := false

	for _, p := range projects {
		if p.ID == projectID {
			found = true
			continue
		}
		updated = append(updated, p)
	}

	if !found {
		return ErrProjectNotFound
	}

	return storage.SaveProjects(updated)
}
