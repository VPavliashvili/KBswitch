package repositories

import "kbswitch/internal/core/switches/models"

type SwitchesRepo interface {
	GetID(brand, name string) (*int, error)
	GetAll() ([]models.SwitchEntity, error)
	GetSingle(int) (*models.SwitchEntity, error)
	AddNew(models.SwitchEntity) (*int, error)
	Remove(int) error
	Update(int, models.SwitchEntity) (*models.SwitchEntity, error)
}
