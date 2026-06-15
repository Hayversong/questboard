package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Hayversong/questboard/internal/model"
)

const (
	EnvDataFile = "QUESTBOARD_DATA_FILE"
	EnvDataDir  = "QUESTBOARD_DATA_DIR"

	defaultDataDir  = "data"
	defaultDataFile = "projects.json"
)

func DataFilePath() string {
	if path := os.Getenv(EnvDataFile); path != "" {
		return path
	}

	dir := os.Getenv(EnvDataDir)
	if dir == "" {
		dir = defaultDataDir
	}

	return filepath.Join(dir, defaultDataFile)
}

func SaveProjectsJSON(
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

	path := DataFilePath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := backupFile(path); err != nil {
		return err
	}

	tempFile, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}

	tempPath := tempFile.Name()

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return err
	}

	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return err
	}

	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return err
	}

	if err := replaceFile(tempPath, path); err != nil {
		os.Remove(tempPath)
		return err
	}

	return nil
}

func LoadProjectsJSON() ([]model.Project, error) {

	data, err := os.ReadFile(
		DataFilePath(),
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

func backupFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath(path), data, 0644)
}

func backupPath(path string) string {
	return path + ".bak"
}

func replaceFile(tempPath string, targetPath string) error {
	if runtime.GOOS == "windows" {
		if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return os.Rename(tempPath, targetPath)
}
