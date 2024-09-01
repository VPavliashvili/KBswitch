package switches

import "kbswitch/internal/core/switches/models"

type Service interface {
	GetAll() ([]models.Switch, error)
	GetSingle(string, string) (*models.Switch, error)
	AddNew(models.SwitchRequestBody) (*int, error)
	Remove(string, string) *models.AppError
	Update(brand, name string, body models.SwitchRequestBody) (*models.Switch, error)
}

type Repo interface {
	GetID(brand, name string) (*int, error)
	GetAll() ([]models.SwitchEntity, error)
	GetSingle(int) (*models.SwitchEntity, error)
	AddNew(models.SwitchEntity) (*int, error)
	Remove(int) error
	Update(int, models.SwitchEntity) (*models.SwitchEntity, error)
}
