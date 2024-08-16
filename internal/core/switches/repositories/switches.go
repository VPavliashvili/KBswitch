package repositories

import "kbswitch/internal/core/switches/models"

type SwitchesRepo interface {
	GetAll() ([]models.SwitchEntity, error)
	GetByID(int) (*models.SwitchEntity, error)
	AddNew(models.SwitchEntity) error
}
