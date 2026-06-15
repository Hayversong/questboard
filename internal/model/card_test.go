package model

import (
	"strings"
	"testing"
	"time"
)

func TestCardCreation(t *testing.T) {
	card := Card{
		ID:          "card-1",
		Title:       "Create inventory",
		Description: "Add item slots",
		Status:      "backlog",
		Rarity:      "rare",
		Deadline:    "2026-06-20",
		Order:       2,
	}

	if card.ID != "card-1" {
		t.Fatalf("expected ID card-1, got %q", card.ID)
	}

	if card.Title == "" {
		t.Fatal("expected title to be set")
	}

	if card.Status != "backlog" {
		t.Fatalf("expected status backlog, got %q", card.Status)
	}

	if card.Order != 2 {
		t.Fatalf("expected order 2, got %d", card.Order)
	}
}

func TestCardZeroValueKeepsRequiredAttributesEmpty(t *testing.T) {
	card := Card{}

	if card.ID != "" {
		t.Fatalf("expected empty ID, got %q", card.ID)
	}

	if card.Title != "" {
		t.Fatalf("expected empty title, got %q", card.Title)
	}
}

func TestCardXPByRarity(t *testing.T) {
	tests := []struct {
		name   string
		rarity string
		want   int
	}{
		{name: "common", rarity: "common", want: 50},
		{name: "rare", rarity: "rare", want: 150},
		{name: "epic", rarity: "epic", want: 300},
		{name: "legendary", rarity: "legendary", want: 600},
		{name: "unknown defaults to common", rarity: "unknown", want: 50},
		{name: "empty defaults to common", rarity: "", want: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := Card{Rarity: tt.rarity}

			if got := card.XP(); got != tt.want {
				t.Fatalf("expected XP %d, got %d", tt.want, got)
			}
		})
	}
}

func TestCardLabels(t *testing.T) {
	tests := []struct {
		name       string
		card       Card
		wantRarity string
		wantStatus string
		wantColor  string
	}{
		{
			name:       "common backlog",
			card:       Card{Rarity: "common", Status: "backlog"},
			wantRarity: "Comum",
			wantStatus: "A Fazer",
			wantColor:  "#555",
		},
		{
			name:       "rare doing",
			card:       Card{Rarity: "rare", Status: "doing"},
			wantRarity: "Rara",
			wantStatus: "Fazendo",
			wantColor:  "#4a90e2",
		},
		{
			name:       "epic done",
			card:       Card{Rarity: "epic", Status: "done"},
			wantRarity: "Épica",
			wantStatus: "Concluído",
			wantColor:  "#8e44ad",
		},
		{
			name:       "legendary",
			card:       Card{Rarity: "legendary", Status: "done"},
			wantRarity: "Lendária",
			wantStatus: "Concluído",
			wantColor:  "#f39c12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.RarityLabel(); got != tt.wantRarity {
				t.Fatalf("expected rarity label %q, got %q", tt.wantRarity, got)
			}

			if got := tt.card.StatusLabel(); got != tt.wantStatus {
				t.Fatalf("expected status label %q, got %q", tt.wantStatus, got)
			}

			if got := tt.card.Color(); got != tt.wantColor {
				t.Fatalf("expected color %q, got %q", tt.wantColor, got)
			}
		})
	}
}

func TestCardIsLate(t *testing.T) {
	past := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	future := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	tests := []struct {
		name string
		card Card
		want bool
	}{
		{name: "empty deadline", card: Card{Status: "backlog"}, want: false},
		{name: "invalid deadline", card: Card{Status: "backlog", Deadline: "invalid"}, want: false},
		{name: "future deadline", card: Card{Status: "backlog", Deadline: future}, want: false},
		{name: "past deadline", card: Card{Status: "backlog", Deadline: past}, want: true},
		{name: "done ignores past deadline", card: Card{Status: "done", Deadline: past}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.IsLate(); got != tt.want {
				t.Fatalf("expected IsLate %v, got %v", tt.want, got)
			}
		})
	}
}

func TestCardDeadlineLabel(t *testing.T) {
	if got := (Card{}).DeadlineLabel(); got != "" {
		t.Fatalf("expected empty label without deadline, got %q", got)
	}

	if got := (Card{Deadline: "invalid"}).DeadlineLabel(); got != "" {
		t.Fatalf("expected empty label with invalid deadline, got %q", got)
	}

	if got := (Card{Status: "done", Deadline: "2026-01-01"}).DeadlineLabel(); got != "Concluída" {
		t.Fatalf("expected done label, got %q", got)
	}

	today := time.Now().Format("2006-01-02")
	if got := (Card{Status: "backlog", Deadline: today}).DeadlineLabel(); got != "Vence hoje" {
		t.Fatalf("expected today label, got %q", got)
	}

	past := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	if got := (Card{Status: "doing", Deadline: past}).DeadlineLabel(); !strings.HasPrefix(got, "Atrasada há ") {
		t.Fatalf("expected late label, got %q", got)
	}
}
