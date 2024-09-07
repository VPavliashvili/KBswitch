package repo

import (
	"context"
	"kbswitch/internal/app"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"
)

func New(config app.DbConfig) switches.Repo {
	return repo{
		user: config.User,
		pass: config.Pass,
		host: config.Host,
		db:   config.Db,
		port: config.Port,
	}
}

type repo struct {
	user string
	pass string
	host string
	db   string
	port int
}

// AddNew implements switches.Repo.
func (r repo) AddNew(context.Context, models.SwitchEntity) (*int, error) {
	panic("unimplemented")
}

// GetAll implements switches.Repo.
func (r repo) GetAll(context.Context) ([]models.SwitchEntity, error) {
	panic("unimplemented")
}

// GetID implements switches.Repo.
func (r repo) GetID(ctx context.Context, brand string, name string) (*int, error) {
	panic("unimplemented")
}

// GetSingle implements switches.Repo.
func (r repo) GetSingle(context.Context, int) (*models.SwitchEntity, error) {
	panic("unimplemented")
}

// Remove implements switches.Repo.
func (r repo) Remove(context.Context, int) error {
	panic("unimplemented")
}

// Update implements switches.Repo.
func (r repo) Update(context.Context, int, models.SwitchEntity) (*models.SwitchEntity, error) {
	panic("unimplemented")
}
