package service

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Hayversong/questboard/internal/model"
	"github.com/Hayversong/questboard/internal/storage"
)

func useTempStorage(t *testing.T, projects []model.Project) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "projects.json")
	t.Setenv(storage.EnvDataFile, path)

	if projects != nil {
		if err := storage.SaveProjects(projects); err != nil {
			t.Fatalf("seed projects: %v", err)
		}
	}

	return path
}

func loadProjects(t *testing.T) []model.Project {
	t.Helper()

	projects, err := storage.LoadProjects()
	if err != nil {
		t.Fatalf("load projects: %v", err)
	}

	return projects
}

func TestCreateProject(t *testing.T) {
	useTempStorage(t, []model.Project{})

	if err := CreateProject("  Dungeon Survivor  "); err != nil {
		t.Fatalf("create project: %v", err)
	}

	projects := loadProjects(t)

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	if projects[0].ID == "" {
		t.Fatal("expected generated project ID")
	}

	if projects[0].Name != "Dungeon Survivor" {
		t.Fatalf("expected trimmed name, got %q", projects[0].Name)
	}
}

func TestCreateProjectRequiresName(t *testing.T) {
	useTempStorage(t, []model.Project{})

	if err := CreateProject("  "); !errors.Is(err, ErrProjectNameRequired) {
		t.Fatalf("expected ErrProjectNameRequired, got %v", err)
	}
}

func TestDashboardStats(t *testing.T) {
	useTempStorage(t, []model.Project{
		{
			ID:   "project-1",
			Name: "Dungeon Survivor",
			Cards: []model.Card{
				{ID: "card-1", Status: "done", Rarity: "rare"},
				{ID: "card-2", Status: "doing", Rarity: "legendary"},
			},
		},
		{
			ID:   "project-2",
			Name: "Platformer",
			Cards: []model.Card{
				{ID: "card-3", Status: "done", Rarity: "common"},
			},
		},
	})

	stats, err := DashboardStats()
	if err != nil {
		t.Fatalf("dashboard stats: %v", err)
	}

	if stats.Projects != 2 || stats.Cards != 3 || stats.Done != 2 || stats.XP != 200 {
		t.Fatalf("unexpected stats: %#v", stats)
	}
}

func TestGetProjectByIDSortsCards(t *testing.T) {
	useTempStorage(t, []model.Project{
		{
			ID:   "project-1",
			Name: "Dungeon Survivor",
			Cards: []model.Card{
				{ID: "card-2", Order: 2},
				{ID: "card-1", Order: 1},
			},
		},
	})

	project, err := GetProjectByID("project-1")
	if err != nil {
		t.Fatalf("get project: %v", err)
	}

	if project.Cards[0].ID != "card-1" {
		t.Fatalf("expected first card sorted by order, got %q", project.Cards[0].ID)
	}
}

func TestGetProjectByIDErrors(t *testing.T) {
	useTempStorage(t, []model.Project{})

	if _, err := GetProjectByID(""); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if _, err := GetProjectByID("missing"); !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestCreateCard(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Name: "Dungeon Survivor"}})

	err := CreateCard("project-1", " Inventory ", " Items ", "epic", "2026-06-20")
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	projects := loadProjects(t)
	card := projects[0].Cards[0]

	if card.ID == "" {
		t.Fatal("expected generated card ID")
	}

	if card.Title != "Inventory" || card.Description != "Items" {
		t.Fatalf("expected trimmed card fields, got %#v", card)
	}

	if card.Status != "backlog" || card.Rarity != "epic" || card.Order != 0 {
		t.Fatalf("unexpected card defaults: %#v", card)
	}

	if len(projects[0].Activities) != 1 {
		t.Fatalf("expected activity on create, got %d", len(projects[0].Activities))
	}
}

func TestCreateCardErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Name: "Dungeon Survivor"}})

	tests := []struct {
		name string
		err  error
		run  func() error
	}{
		{name: "project required", err: ErrProjectIDRequired, run: func() error {
			return CreateCard("", "Title", "", "common", "")
		}},
		{name: "title required", err: ErrCardTitleRequired, run: func() error {
			return CreateCard("project-1", " ", "", "common", "")
		}},
		{name: "invalid rarity", err: ErrInvalidRarity, run: func() error {
			return CreateCard("project-1", "Title", "", "invalid", "")
		}},
		{name: "invalid deadline", err: ErrInvalidDeadline, run: func() error {
			return CreateCard("project-1", "Title", "", "common", "20/06/2026")
		}},
		{name: "project not found", err: ErrProjectNotFound, run: func() error {
			return CreateCard("missing", "Title", "", "common", "")
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestMoveCard(t *testing.T) {
	useTempStorage(t, []model.Project{
		{
			ID:   "project-1",
			Name: "Dungeon Survivor",
			Cards: []model.Card{
				{ID: "card-1", Title: "Inventory", Status: "backlog"},
			},
		},
	})

	if err := MoveCard("project-1", "card-1"); err != nil {
		t.Fatalf("move card to doing: %v", err)
	}

	projects := loadProjects(t)
	if projects[0].Cards[0].Status != "doing" {
		t.Fatalf("expected doing, got %q", projects[0].Cards[0].Status)
	}

	if err := MoveCard("project-1", "card-1"); err != nil {
		t.Fatalf("move card to done: %v", err)
	}

	projects = loadProjects(t)
	if projects[0].Cards[0].Status != "done" {
		t.Fatalf("expected done, got %q", projects[0].Cards[0].Status)
	}

	if len(projects[0].Activities) != 2 {
		t.Fatalf("expected 2 activities, got %d", len(projects[0].Activities))
	}
}

func TestMoveCardErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Name: "Dungeon Survivor"}})

	if err := MoveCard("", "card-1"); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if err := MoveCard("project-1", ""); !errors.Is(err, ErrCardIDRequired) {
		t.Fatalf("expected ErrCardIDRequired, got %v", err)
	}

	if err := MoveCard("project-1", "missing"); !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("expected ErrCardNotFound, got %v", err)
	}

	if err := MoveCard("missing", "card-1"); !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestUpdateCardStatus(t *testing.T) {
	useTempStorage(t, []model.Project{
		{ID: "project-1", Cards: []model.Card{{ID: "card-1", Status: "backlog"}}},
	})

	if err := UpdateCardStatus("project-1", "card-1", "done"); err != nil {
		t.Fatalf("update card status: %v", err)
	}

	projects := loadProjects(t)
	if projects[0].Cards[0].Status != "done" {
		t.Fatalf("expected done, got %q", projects[0].Cards[0].Status)
	}
}

func TestUpdateCardStatusErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Cards: []model.Card{{ID: "card-1"}}}})

	if err := UpdateCardStatus("", "card-1", "done"); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if err := UpdateCardStatus("project-1", "", "done"); !errors.Is(err, ErrCardIDRequired) {
		t.Fatalf("expected ErrCardIDRequired, got %v", err)
	}

	if err := UpdateCardStatus("project-1", "card-1", "invalid"); !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestUpdateCard(t *testing.T) {
	useTempStorage(t, []model.Project{
		{ID: "project-1", Cards: []model.Card{{ID: "card-1", Title: "Old", Rarity: "common"}}},
	})

	err := UpdateCard("project-1", "card-1", " New ", " Desc ", "legendary", "2026-06-20")
	if err != nil {
		t.Fatalf("update card: %v", err)
	}

	projects := loadProjects(t)
	card := projects[0].Cards[0]

	if card.Title != "New" || card.Description != "Desc" || card.Rarity != "legendary" || card.Deadline != "2026-06-20" {
		t.Fatalf("unexpected updated card: %#v", card)
	}
}

func TestUpdateCardErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Cards: []model.Card{{ID: "card-1"}}}})

	if err := UpdateCard("", "card-1", "Title", "", "common", ""); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if err := UpdateCard("project-1", "", "Title", "", "common", ""); !errors.Is(err, ErrCardIDRequired) {
		t.Fatalf("expected ErrCardIDRequired, got %v", err)
	}

	if err := UpdateCard("project-1", "card-1", " ", "", "common", ""); !errors.Is(err, ErrCardTitleRequired) {
		t.Fatalf("expected ErrCardTitleRequired, got %v", err)
	}

	if err := UpdateCard("project-1", "card-1", "Title", "", "invalid", ""); !errors.Is(err, ErrInvalidRarity) {
		t.Fatalf("expected ErrInvalidRarity, got %v", err)
	}
}

func TestDeleteCard(t *testing.T) {
	useTempStorage(t, []model.Project{
		{ID: "project-1", Cards: []model.Card{{ID: "card-1", Title: "Delete me"}}},
	})

	if err := DeleteCard("project-1", "card-1"); err != nil {
		t.Fatalf("delete card: %v", err)
	}

	projects := loadProjects(t)
	if len(projects[0].Cards) != 0 {
		t.Fatalf("expected no cards, got %d", len(projects[0].Cards))
	}

	if len(projects[0].Activities) != 1 {
		t.Fatalf("expected delete activity, got %d", len(projects[0].Activities))
	}
}

func TestDeleteCardErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1"}})

	if err := DeleteCard("", "card-1"); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if err := DeleteCard("project-1", ""); !errors.Is(err, ErrCardIDRequired) {
		t.Fatalf("expected ErrCardIDRequired, got %v", err)
	}

	if err := DeleteCard("project-1", "missing"); !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("expected ErrCardNotFound, got %v", err)
	}
}

func TestReorderCardsUpdatesOrderAndStatus(t *testing.T) {
	useTempStorage(t, []model.Project{
		{ID: "project-1", Cards: []model.Card{
			{ID: "card-1", Status: "backlog", Order: 0},
			{ID: "card-2", Status: "doing", Order: 1},
		}},
	})

	err := ReorderCards("project-1", []CardOrder{
		{ID: "card-1", Order: 1, Status: "done"},
		{ID: "card-2", Order: 0, Status: "doing"},
	})
	if err != nil {
		t.Fatalf("reorder cards: %v", err)
	}

	projects := loadProjects(t)

	if projects[0].Cards[0].Order != 1 || projects[0].Cards[0].Status != "done" {
		t.Fatalf("unexpected first card: %#v", projects[0].Cards[0])
	}

	if projects[0].Cards[1].Order != 0 || projects[0].Cards[1].Status != "doing" {
		t.Fatalf("unexpected second card: %#v", projects[0].Cards[1])
	}
}

func TestReorderCardsErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Cards: []model.Card{{ID: "card-1"}}}})

	tests := []struct {
		name string
		err  error
		run  func() error
	}{
		{name: "project required", err: ErrProjectIDRequired, run: func() error {
			return ReorderCards("", []CardOrder{{ID: "card-1"}})
		}},
		{name: "cards required", err: ErrCardsListRequired, run: func() error {
			return ReorderCards("project-1", nil)
		}},
		{name: "card id required", err: ErrCardIDRequired, run: func() error {
			return ReorderCards("project-1", []CardOrder{{ID: ""}})
		}},
		{name: "invalid order", err: ErrInvalidOrder, run: func() error {
			return ReorderCards("project-1", []CardOrder{{ID: "card-1", Order: -1}})
		}},
		{name: "invalid status", err: ErrInvalidStatus, run: func() error {
			return ReorderCards("project-1", []CardOrder{{ID: "card-1", Status: "invalid"}})
		}},
		{name: "project not found", err: ErrProjectNotFound, run: func() error {
			return ReorderCards("missing", []CardOrder{{ID: "card-1"}})
		}},
		{name: "card not found", err: ErrCardNotFound, run: func() error {
			return ReorderCards("project-1", []CardOrder{{ID: "missing"}})
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestRenameProject(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Name: "Old"}})

	if err := RenameProject("project-1", " New Name "); err != nil {
		t.Fatalf("rename project: %v", err)
	}

	projects := loadProjects(t)
	if projects[0].Name != "New Name" {
		t.Fatalf("expected renamed project, got %q", projects[0].Name)
	}
}

func TestRenameProjectErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1", Name: "Old"}})

	if err := RenameProject("", "Name"); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if err := RenameProject("project-1", " "); !errors.Is(err, ErrProjectNameRequired) {
		t.Fatalf("expected ErrProjectNameRequired, got %v", err)
	}

	if err := RenameProject("missing", "Name"); !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestDeleteProject(t *testing.T) {
	useTempStorage(t, []model.Project{
		{ID: "project-1", Name: "Delete"},
		{ID: "project-2", Name: "Keep"},
	})

	if err := DeleteProject("project-1"); err != nil {
		t.Fatalf("delete project: %v", err)
	}

	projects := loadProjects(t)
	if len(projects) != 1 || projects[0].ID != "project-2" {
		t.Fatalf("unexpected projects after delete: %#v", projects)
	}
}

func TestDeleteProjectErrors(t *testing.T) {
	useTempStorage(t, []model.Project{{ID: "project-1"}})

	if err := DeleteProject(""); !errors.Is(err, ErrProjectIDRequired) {
		t.Fatalf("expected ErrProjectIDRequired, got %v", err)
	}

	if err := DeleteProject("missing"); !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}
