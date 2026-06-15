package storage

import (
	"errors"
	"os"

	"github.com/Hayversong/questboard/internal/model"
)

var ErrSQLiteDatabaseExists = errors.New("banco SQLite ja existe")

type JSONToSQLiteMigration struct {
	JSONPath   string
	SQLitePath string
	Overwrite  bool
}

type MigrationSummary struct {
	Projects   int
	Cards      int
	Activities int
}

func MigrateJSONToSQLite(options JSONToSQLiteMigration) (MigrationSummary, error) {
	jsonPath := options.JSONPath
	if jsonPath == "" {
		jsonPath = DataFilePath()
	}

	sqlitePath := options.SQLitePath
	if sqlitePath == "" {
		sqlitePath = SQLiteFilePath()
	}

	if _, err := os.Stat(jsonPath); err != nil {
		return MigrationSummary{}, err
	}

	if !options.Overwrite {
		if _, err := os.Stat(sqlitePath); err == nil {
			return MigrationSummary{}, ErrSQLiteDatabaseExists
		} else if !os.IsNotExist(err) {
			return MigrationSummary{}, err
		}
	}

	projects, err := LoadProjectsJSONFromPath(jsonPath)
	if err != nil {
		return MigrationSummary{}, err
	}

	if err := SaveProjectsSQLiteToPath(projects, sqlitePath); err != nil {
		return MigrationSummary{}, err
	}

	return migrationSummary(projects), nil
}

func migrationSummary(projects []model.Project) MigrationSummary {
	summary := MigrationSummary{
		Projects: len(projects),
	}

	for _, project := range projects {
		summary.Cards += len(project.Cards)
		summary.Activities += len(project.Activities)
	}

	return summary
}
