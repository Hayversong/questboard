package model

import (
	"fmt"
	"time"
)

type Card struct {
	ID          string
	Title       string
	Description string
	Status      string
	Rarity      string
	Deadline    string

	Order int
}

func (c Card) IsLate() bool {

	if c.Deadline == "" {
		return false
	}

	deadline, err := time.Parse(
		"2006-01-02",
		c.Deadline,
	)

	if err != nil {
		return false
	}

	return time.Now().After(deadline) &&
		c.Status != "done"
}

func (c Card) DeadlineLabel() string {

	if c.Deadline == "" {
		return ""
	}

	deadline, err := time.Parse(
		"2006-01-02",
		c.Deadline,
	)

	if err != nil {
		return ""
	}

	now := time.Now()

	days := int(
		deadline.Sub(now).Hours() / 24,
	)

	if c.Status == "done" {
		return "Concluída"
	}

	if days < 0 {

		return fmt.Sprintf(
			"Atrasada há %d dias",
			-days,
		)
	}

	if days == 0 {
		return "Vence hoje"
	}

	if days == 1 {
		return "Vence amanhã"
	}

	return fmt.Sprintf(
		"Vence em %d dias",
		days,
	)
}

func (c Card) XP() int {

	switch c.Rarity {

	case "rare":
		return 150

	case "epic":
		return 300

	case "legendary":
		return 600

	default:
		return 50
	}
}

func (c Card) RarityLabel() string {
	switch c.Rarity {
	case "rare":
		return "Rara"
	case "epic":
		return "\u00c9pica"
	case "legendary":
		return "Lend\u00e1ria"
	default:
		return "Comum"
	}
}

func (c Card) StatusLabel() string {
	switch c.Status {
	case "doing":
		return "Fazendo"
	case "done":
		return "Conclu\u00eddo"
	default:
		return "A Fazer"
	}
}

func (c Card) Color() string {

	switch c.Rarity {

	case "rare":
		return "#4a90e2"

	case "epic":
		return "#8e44ad"

	case "legendary":
		return "#f39c12"

	default:
		return "#555"
	}
}
