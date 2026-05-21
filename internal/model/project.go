package model

type Project struct {
	ID    string
	Name  string
	Cards []Card
}

func (p Project) Progress() int {

	total := len(p.Cards)

	if total == 0 {
		return 0
	}

	done := 0

	for _, card := range p.Cards {

		if card.Status == "done" {
			done++
		}
	}

	return (done * 100) / total
}
