package storage

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Hayversong/questboard/internal/model"
)

func TestMigrateJSONToSQLite(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "projects.json")
	sqlitePath := filepath.Join(dir, "questboard.db")

	t.Setenv(EnvDataFile, jsonPath)

	projects := []model.Project{
		{
			ID:   "project-1",
			Name: "Dungeon Survivor",
			Cards: []model.Card{
				{
					ID:          "card-1",
					Title:       "Inventory",
					Description: "Items",
					Status:      "done",
					Rarity:      "epic",
					Deadline:    "2026-06-20",
					Order:       1,
				},
			},
			Activities: []model.Activity{
				{Message: "Quest criada", Time: "15/06 10:00"},
			},
		},
	}

	if err := SaveProjectsJSON(projects); err != nil {
		t.Fatalf("save json projects: %v", err)
	}

	summary, err := MigrateJSONToSQLite(JSONToSQLiteMigration{
		JSONPath:   jsonPath,
		SQLitePath: sqlitePath,
	})
	if err != nil {
		t.Fatalf("migrate json to sqlite: %v", err)
	}

	if summary.Projects != 1 || summary.Cards != 1 || summary.Activities != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}

	t.Setenv(EnvStorageBackend, StorageSQLite)
	t.Setenv(EnvSQLiteFile, sqlitePath)

	got, err := LoadProjects()
	if err != nil {
		t.Fatalf("load migrated sqlite projects: %v", err)
	}

	if len(got) != 1 || got[0].ID != "project-1" {
		t.Fatalf("unexpected migrated projects: %#v", got)
	}

	if len(got[0].Cards) != 1 || got[0].Cards[0].ID != "card-1" {
		t.Fatalf("unexpected migrated cards: %#v", got[0].Cards)
	}

	if len(got[0].Activities) != 1 || got[0].Activities[0].Message != "Quest criada" {
		t.Fatalf("unexpected migrated activities: %#v", got[0].Activities)
	}
}

func TestMigrateJSONToSQLiteDoesNotOverwriteByDefault(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "projects.json")
	sqlitePath := filepath.Join(dir, "questboard.db")

	t.Setenv(EnvDataFile, jsonPath)

	if err := SaveProjectsJSON([]model.Project{{ID: "project-1", Name: "Game"}}); err != nil {
		t.Fatalf("save json projects: %v", err)
	}

	if _, err := MigrateJSONToSQLite(JSONToSQLiteMigration{
		JSONPath:   jsonPath,
		SQLitePath: sqlitePath,
	}); err != nil {
		t.Fatalf("initial migration: %v", err)
	}

	_, err := MigrateJSONToSQLite(JSONToSQLiteMigration{
		JSONPath:   jsonPath,
		SQLitePath: sqlitePath,
	})
	if !errors.Is(err, ErrSQLiteDatabaseExists) {
		t.Fatalf("expected ErrSQLiteDatabaseExists, got %v", err)
	}
}

func TestMigrateJSONToSQLiteOverwrite(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "projects.json")
	sqlitePath := filepath.Join(dir, "questboard.db")

	t.Setenv(EnvDataFile, jsonPath)

	if err := SaveProjectsJSON([]model.Project{{ID: "project-1", Name: "First"}}); err != nil {
		t.Fatalf("save first json projects: %v", err)
	}

	if _, err := MigrateJSONToSQLite(JSONToSQLiteMigration{
		JSONPath:   jsonPath,
		SQLitePath: sqlitePath,
	}); err != nil {
		t.Fatalf("initial migration: %v", err)
	}

	if err := SaveProjectsJSON([]model.Project{{ID: "project-2", Name: "Second"}}); err != nil {
		t.Fatalf("save second json projects: %v", err)
	}

	if _, err := MigrateJSONToSQLite(JSONToSQLiteMigration{
		JSONPath:   jsonPath,
		SQLitePath: sqlitePath,
		Overwrite:  true,
	}); err != nil {
		t.Fatalf("overwrite migration: %v", err)
	}

	t.Setenv(EnvStorageBackend, StorageSQLite)
	t.Setenv(EnvSQLiteFile, sqlitePath)

	got, err := LoadProjects()
	if err != nil {
		t.Fatalf("load overwritten sqlite projects: %v", err)
	}

	if len(got) != 1 || got[0].ID != "project-2" {
		t.Fatalf("unexpected overwritten projects: %#v", got)
	}
}
