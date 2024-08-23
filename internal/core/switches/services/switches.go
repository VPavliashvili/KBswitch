package services

import "kbswitch/internal/core/switches/models"

type SwitchesService interface {
	GetAll() ([]models.Switch, error)
	GetSingle(string, string) (*models.Switch, error)
	AddNew(models.SwitchRequestBody) (*int, error)
	Remove(string, string) error
	Update(models.SwitchRequestBody) (*models.Switch, error)
}
