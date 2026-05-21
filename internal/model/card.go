package model

type Card struct {
	ID          string
	Title       string
	Description string
	Status      string
	Rarity      string
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
