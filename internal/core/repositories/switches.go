package repositories

import "kbswitch/internal/core/models"

type SwitchesRepo interface {
	GetAll() ([]models.SwitchEntity, error)
	GetByID(int) (*models.SwitchEntity, error)
}
