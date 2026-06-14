package service

import (
	"strings"

	"github.com/Hayversong/questboard/internal/storage"
)

func RenameProject(projectID string, name string) error {

	if projectID == "" {
		return ErrProjectIDRequired
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return ErrProjectNameRequired
	}

	projects, err := storage.LoadProjects()
	if err != nil {
		return err
	}

	for i := range projects {

		if projects[i].ID == projectID {

			projects[i].Name = name

			return storage.SaveProjects(projects)
		}
	}

	return ErrProjectNotFound
}
