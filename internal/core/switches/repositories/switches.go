package repositories

import "kbswitch/internal/core/switches/models"

type SwitchesRepo interface {
	GetAll() ([]models.SwitchEntity, error)
	GetSingle(string, string) (*models.SwitchEntity, error)
	AddNew(models.SwitchEntity) error
	Exists(string, string) (bool, error)
}
