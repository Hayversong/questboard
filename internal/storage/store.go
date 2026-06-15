package storage

import (
	"os"
	"strings"

	"github.com/Hayversong/questboard/internal/model"
)

const (
	EnvStorageBackend = "QUESTBOARD_STORAGE"
	StorageJSON       = "json"
	StorageSQLite     = "sqlite"
)

func LoadProjects() ([]model.Project, error) {
	if useSQLite() {
		return LoadProjectsSQLite()
	}

	return LoadProjectsJSON()
}

func SaveProjects(projects []model.Project) error {
	if useSQLite() {
		return SaveProjectsSQLite(projects)
	}

	return SaveProjectsJSON(projects)
}

func useSQLite() bool {
	backend := strings.ToLower(strings.TrimSpace(os.Getenv(EnvStorageBackend)))

	return backend == StorageSQLite || os.Getenv(EnvSQLiteFile) != ""
}
