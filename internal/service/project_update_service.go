package service

import (
	"errors"

	"github.com/Hayversong/questboard/internal/storage"
)

func RenameProject(
	projectID string,
	name string,
) error {

	projects, err := storage.LoadProjects()

	if err != nil {
		return err
	}

	for i := range projects {

		if projects[i].ID == projectID {

			projects[i].Name = name

			return storage.SaveProjects(
				projects,
			)
		}
	}

	return errors.New(
		"projeto não encontrado",
	)
}
