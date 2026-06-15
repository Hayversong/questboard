package model

import "testing"

func TestProjectCreation(t *testing.T) {
	project := Project{
		ID:   "project-1",
		Name: "Dungeon Survivor",
		Cards: []Card{
			{ID: "card-1", Title: "Movement", Status: "done"},
		},
		Activities: []Activity{
			{Message: "Quest criada", Time: "15/06 10:00"},
		},
	}

	if project.ID != "project-1" {
		t.Fatalf("expected ID project-1, got %q", project.ID)
	}

	if project.Name == "" {
		t.Fatal("expected project name to be set")
	}

	if len(project.Cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(project.Cards))
	}

	if len(project.Activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(project.Activities))
	}
}

func TestProjectZeroValueKeepsRequiredAttributesEmpty(t *testing.T) {
	project := Project{}

	if project.ID != "" {
		t.Fatalf("expected empty ID, got %q", project.ID)
	}

	if project.Name != "" {
		t.Fatalf("expected empty name, got %q", project.Name)
	}
}

func TestProjectProgress(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		want    int
	}{
		{name: "without cards", project: Project{}, want: 0},
		{
			name: "partial progress",
			project: Project{Cards: []Card{
				{Status: "backlog"},
				{Status: "doing"},
				{Status: "done"},
			}},
			want: 33,
		},
		{
			name: "completed",
			project: Project{Cards: []Card{
				{Status: "done"},
				{Status: "done"},
			}},
			want: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.project.Progress(); got != tt.want {
				t.Fatalf("expected progress %d, got %d", tt.want, got)
			}
		})
	}
}

func TestProjectXPLevelAndRank(t *testing.T) {
	tests := []struct {
		name      string
		project   Project
		wantXP    int
		wantLevel int
		wantRank  string
	}{
		{name: "empty project", project: Project{}, wantXP: 0, wantLevel: 1, wantRank: "Novato"},
		{
			name: "only done cards count XP",
			project: Project{Cards: []Card{
				{Status: "done", Rarity: "rare"},
				{Status: "doing", Rarity: "legendary"},
			}},
			wantXP:    150,
			wantLevel: 2,
			wantRank:  "Novato",
		},
		{
			name: "adventurer rank",
			project: Project{Cards: []Card{
				{Status: "done", Rarity: "legendary"},
			}},
			wantXP:    600,
			wantLevel: 7,
			wantRank:  "Herói",
		},
		{
			name: "legend rank",
			project: Project{Cards: []Card{
				{Status: "done", Rarity: "legendary"},
				{Status: "done", Rarity: "legendary"},
			}},
			wantXP:    1200,
			wantLevel: 13,
			wantRank:  "Lenda",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.project.XP(); got != tt.wantXP {
				t.Fatalf("expected XP %d, got %d", tt.wantXP, got)
			}

			if got := tt.project.Level(); got != tt.wantLevel {
				t.Fatalf("expected level %d, got %d", tt.wantLevel, got)
			}

			if got := tt.project.Rank(); got != tt.wantRank {
				t.Fatalf("expected rank %q, got %q", tt.wantRank, got)
			}
		})
	}
}

func TestProjectAchievements(t *testing.T) {
	tests := []struct {
		name    string
		project Project
		want    []string
	}{
		{name: "none", project: Project{}, want: nil},
		{
			name: "first quest",
			project: Project{Cards: []Card{
				{Status: "done", Rarity: "common"},
			}},
			want: []string{"Primeira Quest"},
		},
		{
			name: "five done quests",
			project: Project{Cards: []Card{
				{Status: "done", Rarity: "common"},
				{Status: "done", Rarity: "common"},
				{Status: "done", Rarity: "common"},
				{Status: "done", Rarity: "common"},
				{Status: "done", Rarity: "common"},
			}},
			want: []string{"Primeira Quest", "Aventureiro"},
		},
		{
			name: "master project",
			project: Project{Cards: []Card{
				{Status: "done", Rarity: "legendary"},
				{Status: "done", Rarity: "legendary"},
			}},
			want: []string{"Primeira Quest", "Mestre do Projeto"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.project.Achievements()

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d achievements, got %d", len(tt.want), len(got))
			}

			for i, want := range tt.want {
				if got[i].Title != want {
					t.Fatalf("expected achievement %q at index %d, got %q", want, i, got[i].Title)
				}
			}
		})
	}
}

func TestSimpleModelCreation(t *testing.T) {
	activity := Activity{
		Message: "Quest concluída",
		Time:    "15/06 10:00",
	}

	if activity.Message == "" || activity.Time == "" {
		t.Fatalf("expected activity fields to be set: %#v", activity)
	}

	achievement := Achievement{Title: "Primeira Quest"}
	if achievement.Title == "" {
		t.Fatal("expected achievement title to be set")
	}

	dashboard := Dashboard{Projects: 1, Cards: 2, Done: 1, XP: 50}
	if dashboard.Projects != 1 || dashboard.Cards != 2 || dashboard.Done != 1 || dashboard.XP != 50 {
		t.Fatalf("unexpected dashboard values: %#v", dashboard)
	}
}
