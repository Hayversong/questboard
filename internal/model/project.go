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

func (p Project) XP() int {

	totalXP := 0

	for _, card := range p.Cards {

		if card.Status == "done" {
			totalXP += card.XP()
		}
	}

	return totalXP
}

func (p Project) Level() int {

	return (p.XP() / 100) + 1
}

func (p Project) Achievements() []Achievement {

	var achievements []Achievement

	done := 0

	for _, card := range p.Cards {

		if card.Status == "done" {
			done++
		}
	}

	if done >= 1 {
		achievements = append(
			achievements,
			Achievement{
				Title: "Primeira Quest",
			},
		)
	}

	if done >= 5 {
		achievements = append(
			achievements,
			Achievement{
				Title: "Aventureiro",
			},
		)
	}

	if p.Level() >= 10 {
		achievements = append(
			achievements,
			Achievement{
				Title: "Mestre do Projeto",
			},
		)
	}

	if done >= 20 {
		achievements = append(
			achievements,
			Achievement{
				Title: "Lenda",
			},
		)
	}

	return achievements
}
