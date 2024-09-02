package repo

import (
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
	port uint16
}

// AddNew implements switches.Repo.
func (r repo) AddNew(models.SwitchEntity) (*int, error) {
	panic("unimplemented")
}

// GetAll implements switches.Repo.
func (r repo) GetAll() ([]models.SwitchEntity, error) {
	panic("unimplemented")
}

// GetID implements switches.Repo.
func (r repo) GetID(brand string, name string) (*int, error) {
	panic("unimplemented")
}

// GetSingle implements switches.Repo.
func (r repo) GetSingle(int) (*models.SwitchEntity, error) {
	panic("unimplemented")
}

// Remove implements switches.Repo.
func (r repo) Remove(int) error {
	panic("unimplemented")
}

// Update implements switches.Repo.
func (r repo) Update(int, models.SwitchEntity) (*models.SwitchEntity, error) {
	panic("unimplemented")
}
