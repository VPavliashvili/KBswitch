package switches

import (
	"context"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/switches/models"
)

type Service interface {
	GetAll(context.Context) ([]models.Switch, *common.AppError)
	GetSingle(context.Context, string, string) (*models.Switch, *common.AppError)
	AddNew(context.Context, models.SwitchRequestBody) (*int, *common.AppError)
	Remove(context.Context, string, string) *common.AppError
	Update(ctx context.Context, brand, name string, body models.SwitchRequestBody) (*models.Switch, *common.AppError)
}

type Repo interface {
	GetID(ctx context.Context, brand, name string) (*int, error)
	GetAll(context.Context) ([]models.SwitchEntity, error)
	GetSingle(context.Context, int) (*models.SwitchEntity, error)
	AddNew(context.Context, models.SwitchEntity) (*int, error)
	Remove(context.Context, int) error
	Update(context.Context, int, models.SwitchEntity) (*models.SwitchEntity, error)
}
