package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Hayversong/questboard/internal/model"
)

func TestDataFilePathUsesExplicitFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "custom.json")
	t.Setenv(EnvDataFile, path)
	t.Setenv(EnvDataDir, filepath.Join(t.TempDir(), "ignored"))

	if got := DataFilePath(); got != path {
		t.Fatalf("expected explicit data file %q, got %q", path, got)
	}
}

func TestDataFilePathUsesConfiguredDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(EnvDataFile, "")
	t.Setenv(EnvDataDir, dir)

	want := filepath.Join(dir, defaultDataFile)
	if got := DataFilePath(); got != want {
		t.Fatalf("expected configured data path %q, got %q", want, got)
	}
}

func TestLoadProjectsReturnsEmptyWhenFileDoesNotExist(t *testing.T) {
	t.Setenv(EnvDataFile, filepath.Join(t.TempDir(), "missing.json"))

	projects, err := LoadProjects()
	if err != nil {
		t.Fatalf("load missing projects: %v", err)
	}

	if len(projects) != 0 {
		t.Fatalf("expected empty project list, got %d", len(projects))
	}
}

func TestLoadProjectsReturnsEmptyWhenFileIsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "projects.json")
	t.Setenv(EnvDataFile, path)

	if err := os.WriteFile(path, nil, 0644); err != nil {
		t.Fatalf("write empty file: %v", err)
	}

	projects, err := LoadProjects()
	if err != nil {
		t.Fatalf("load empty projects: %v", err)
	}

	if len(projects) != 0 {
		t.Fatalf("expected empty project list, got %d", len(projects))
	}
}

func TestSaveProjectsWritesAtomicallyAndLoadsData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "projects.json")
	t.Setenv(EnvDataFile, path)

	want := []model.Project{
		{ID: "project-1", Name: "Dungeon Survivor"},
	}

	if err := SaveProjects(want); err != nil {
		t.Fatalf("save projects: %v", err)
	}

	got, err := LoadProjects()
	if err != nil {
		t.Fatalf("load projects: %v", err)
	}

	if len(got) != 1 || got[0].ID != "project-1" || got[0].Name != "Dungeon Survivor" {
		t.Fatalf("unexpected loaded projects: %#v", got)
	}

	matches, err := filepath.Glob(filepath.Join(filepath.Dir(path), "*.tmp"))
	if err != nil {
		t.Fatalf("glob temp files: %v", err)
	}

	if len(matches) != 0 {
		t.Fatalf("expected no temp files left behind, got %#v", matches)
	}
}

func TestSaveProjectsCreatesSingleBackupBeforeOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "projects.json")
	t.Setenv(EnvDataFile, path)

	first := []model.Project{{ID: "project-1", Name: "First"}}
	second := []model.Project{{ID: "project-2", Name: "Second"}}

	if err := SaveProjects(first); err != nil {
		t.Fatalf("save first projects: %v", err)
	}

	if _, err := os.Stat(backupPath(path)); !os.IsNotExist(err) {
		t.Fatalf("expected no backup on first write, got err %v", err)
	}

	if err := SaveProjects(second); err != nil {
		t.Fatalf("save second projects: %v", err)
	}

	current, err := LoadProjects()
	if err != nil {
		t.Fatalf("load current projects: %v", err)
	}

	if len(current) != 1 || current[0].ID != "project-2" {
		t.Fatalf("unexpected current projects: %#v", current)
	}

	t.Setenv(EnvDataFile, backupPath(path))

	backup, err := LoadProjects()
	if err != nil {
		t.Fatalf("load backup projects: %v", err)
	}

	if len(backup) != 1 || backup[0].ID != "project-1" {
		t.Fatalf("unexpected backup projects: %#v", backup)
	}
}

func TestLoadProjectsReturnsInvalidJSONError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "projects.json")
	t.Setenv(EnvDataFile, path)

	if err := os.WriteFile(path, []byte("{invalid"), 0644); err != nil {
		t.Fatalf("write invalid json: %v", err)
	}

	if _, err := LoadProjects(); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}
