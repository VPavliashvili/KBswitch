package switches

import (
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/switches/models"
)

type Service interface {
	GetAll() ([]models.Switch, *common.AppError)
	GetSingle(string, string) (*models.Switch, *common.AppError)
	AddNew(models.SwitchRequestBody) (*int, *common.AppError)
	Remove(string, string) *common.AppError
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
