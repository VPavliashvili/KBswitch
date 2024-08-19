package services

import "kbswitch/internal/core/switches/models"

type SwitchesService interface {
	GetAll() ([]models.Switch, error)
	GetByID(int) (*models.Switch, error)
	AddNew(models.SwitchRequestBody) (*int, error)
	Remove(string, string) error
}
