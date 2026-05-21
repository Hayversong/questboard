package storage

import (
	"encoding/json"
	"os"

	"github.com/Hayversong/questboard/internal/model"
)

func SaveProjects(
	projects []model.Project,
) error {

	data, err := json.MarshalIndent(
		projects,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		"data/projects.json",
		data,
		0644,
	)
}

func LoadProjects() ([]model.Project, error) {

	data, err := os.ReadFile(
		"data/projects.json",
	)

	if err != nil {

		if os.IsNotExist(err) {
			return []model.Project{}, nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return []model.Project{}, nil
	}

	var projects []model.Project

	err = json.Unmarshal(
		data,
		&projects,
	)

	if err != nil {
		return nil, err
	}

	return projects, nil
}
