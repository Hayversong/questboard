package storage

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/Hayversong/questboard/internal/model"
)

func useTempSQLite(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "questboard.db")

	t.Setenv(EnvStorageBackend, StorageSQLite)
	t.Setenv(EnvSQLiteFile, path)

	return path
}

func TestSQLiteSaveAndLoadProjects(t *testing.T) {
	useTempSQLite(t)

	projects := []model.Project{
		{
			ID:   "project-1",
			Name: "Dungeon Survivor",
			Cards: []model.Card{
				{
					ID:          "card-1",
					Title:       "Inventory",
					Description: "Add item slots",
					Status:      "doing",
					Rarity:      "epic",
					Deadline:    "2026-06-20",
					Order:       2,
				},
			},
			Activities: []model.Activity{
				{Message: "Quest criada: Inventory", Time: "15/06 10:00"},
			},
		},
	}

	if err := SaveProjects(projects); err != nil {
		t.Fatalf("save projects with sqlite: %v", err)
	}

	got, err := LoadProjects()
	if err != nil {
		t.Fatalf("load projects with sqlite: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 project, got %d", len(got))
	}

	if got[0].ID != "project-1" || got[0].Name != "Dungeon Survivor" {
		t.Fatalf("unexpected project: %#v", got[0])
	}

	if len(got[0].Cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(got[0].Cards))
	}

	card := got[0].Cards[0]
	if card.ID != "card-1" || card.Status != "doing" || card.Rarity != "epic" || card.Order != 2 {
		t.Fatalf("unexpected card: %#v", card)
	}

	if len(got[0].Activities) != 1 || got[0].Activities[0].Message != "Quest criada: Inventory" {
		t.Fatalf("unexpected activities: %#v", got[0].Activities)
	}
}

func TestSQLiteProjectRepositoryCRUD(t *testing.T) {
	useTempSQLite(t)

	db, err := OpenSQLite()
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	repository := NewProjectRepositorySQLite(db)

	project := projectRecord{
		ID:          "project-1",
		Name:        "Old",
		Description: "Prototype",
	}

	if err := repository.Create(project); err != nil {
		t.Fatalf("create project: %v", err)
	}

	found, err := repository.FindByID("project-1")
	if err != nil {
		t.Fatalf("find project: %v", err)
	}

	if found.Name != "Old" || found.Description != "Prototype" {
		t.Fatalf("unexpected project found: %#v", found)
	}

	project.Name = "New"
	if err := repository.Update(project); err != nil {
		t.Fatalf("update project: %v", err)
	}

	found, err = repository.FindByID("project-1")
	if err != nil {
		t.Fatalf("find updated project: %v", err)
	}

	if found.Name != "New" {
		t.Fatalf("expected updated name, got %q", found.Name)
	}

	if err := repository.Delete("project-1"); err != nil {
		t.Fatalf("delete project: %v", err)
	}

	if _, err := repository.FindByID("project-1"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestSQLiteCardRepositoryCRUD(t *testing.T) {
	useTempSQLite(t)

	db, err := OpenSQLite()
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	projects := NewProjectRepositorySQLite(db)
	cards := NewCardRepositorySQLite(db)

	if err := projects.Create(projectRecord{ID: "project-1", Name: "Game"}); err != nil {
		t.Fatalf("create project: %v", err)
	}

	card := cardRecord{
		ID:          "card-1",
		ProjectID:   "project-1",
		Title:       "Old",
		Description: "First",
		Status:      "backlog",
		Rarity:      "common",
		Order:       1,
	}

	if err := cards.Create(card); err != nil {
		t.Fatalf("create card: %v", err)
	}

	card.Title = "New"
	card.Status = "done"
	card.Rarity = "legendary"
	card.Order = 3

	if err := cards.Update(card); err != nil {
		t.Fatalf("update card: %v", err)
	}

	found, err := cards.FindByID("card-1")
	if err != nil {
		t.Fatalf("find card: %v", err)
	}

	if found.Title != "New" || found.Status != "done" || found.Rarity != "legendary" || found.Order != 3 {
		t.Fatalf("unexpected card: %#v", found)
	}

	list, err := cards.ListByProject("project-1")
	if err != nil {
		t.Fatalf("list cards: %v", err)
	}

	if len(list) != 1 || list[0].ID != "card-1" {
		t.Fatalf("unexpected card list: %#v", list)
	}

	if err := cards.Delete("card-1"); err != nil {
		t.Fatalf("delete card: %v", err)
	}
}

func TestSQLiteActivityRepositoryCRUD(t *testing.T) {
	useTempSQLite(t)

	db, err := OpenSQLite()
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	projects := NewProjectRepositorySQLite(db)
	cards := NewCardRepositorySQLite(db)
	activities := NewActivityRepositorySQLite(db)

	if err := projects.Create(projectRecord{ID: "project-1", Name: "Game"}); err != nil {
		t.Fatalf("create project: %v", err)
	}

	if err := cards.Create(cardRecord{
		ID:        "card-1",
		ProjectID: "project-1",
		Title:     "Card",
		Status:    "backlog",
		Rarity:    "common",
	}); err != nil {
		t.Fatalf("create card: %v", err)
	}

	activity := activityRecord{
		ID:          "activity-1",
		ProjectID:   "project-1",
		CardID:      "card-1",
		Description: "Quest criada",
	}

	if err := activities.Create(activity); err != nil {
		t.Fatalf("create activity: %v", err)
	}

	activity.Description = "Quest atualizada"
	if err := activities.Update(activity); err != nil {
		t.Fatalf("update activity: %v", err)
	}

	found, err := activities.FindByID("activity-1")
	if err != nil {
		t.Fatalf("find activity: %v", err)
	}

	if found.Description != "Quest atualizada" || found.CardID != "card-1" {
		t.Fatalf("unexpected activity: %#v", found)
	}

	list, err := activities.ListByProject("project-1")
	if err != nil {
		t.Fatalf("list activities: %v", err)
	}

	if len(list) != 1 || list[0].ID != "activity-1" {
		t.Fatalf("unexpected activity list: %#v", list)
	}

	if err := activities.Delete("activity-1"); err != nil {
		t.Fatalf("delete activity: %v", err)
	}
}

func TestSQLiteForeignKeyValidation(t *testing.T) {
	useTempSQLite(t)

	db, err := OpenSQLite()
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	cards := NewCardRepositorySQLite(db)

	err = cards.Create(cardRecord{
		ID:        "card-1",
		ProjectID: "missing",
		Title:     "Invalid",
		Status:    "backlog",
		Rarity:    "common",
	})
	if err == nil {
		t.Fatal("expected foreign key error")
	}
}
