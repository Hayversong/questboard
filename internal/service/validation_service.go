package service

import (
	"errors"
	"time"
)

var (
	ErrProjectIDRequired   = errors.New("id do projeto obrigatório")
	ErrCardIDRequired      = errors.New("id do card obrigatório")
	ErrProjectNameRequired = errors.New("nome do projeto obrigatório")
	ErrCardTitleRequired   = errors.New("título obrigatório")
	ErrInvalidRarity       = errors.New("raridade inválida")
	ErrInvalidStatus       = errors.New("status inválido")
	ErrInvalidDeadline     = errors.New("data limite inválida")
	ErrProjectNotFound     = errors.New("projeto não encontrado")
	ErrCardNotFound        = errors.New("card não encontrado")
	ErrCardsListRequired   = errors.New("lista de cards obrigatória")
	ErrInvalidOrder        = errors.New("ordem inválida")
)

func isValidRarity(rarity string) bool {
	switch rarity {
	case "common", "rare", "epic", "legendary":
		return true
	default:
		return false
	}
}

func isValidStatus(status string) bool {
	switch status {
	case "backlog", "doing", "done":
		return true
	default:
		return false
	}
}

func isValidDeadline(deadline string) bool {
	if deadline == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", deadline)
	return err == nil
}
