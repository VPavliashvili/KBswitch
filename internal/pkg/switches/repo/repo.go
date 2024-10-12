package repo

import (
	"context"
	"fmt"
	"kbswitch/internal/app"
	"kbswitch/internal/core/common/database"
	"kbswitch/internal/core/common/logging"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"
)

func New(logger logging.Logger, pool database.DBPool) switches.Repo {
	return repo{
		pool:   pool,
		logger: logger,
	}
}

func NewObsolete(cfg app.DbConfig) switches.Repo {
	return repo{
		cfg: cfg,
	}
}

type repo struct {
	logger logging.Logger
	pool   database.DBPool
	cfg    app.DbConfig
}

// AddNew implements switches.Repo.
func (r repo) AddNew(context.Context, models.SwitchEntity) (*int, error) {
	panic("unimplemented")
}

// GetAll implements switches.Repo.
func (r repo) GetAll(ctx context.Context) ([]models.SwitchEntity, error) {
	result := make([]models.SwitchEntity, 0)
	query := `SELECT * FROM public.switches`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.logger.LogError(fmt.Sprintf("query error: %s", err.Error()))
		return []models.SwitchEntity{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.SwitchEntity
		rows.Scan(&r.ID, &r.Manufacturer, &r.ActuationType, &r.Lifespan,
			&r.Model, &r.Image, &r.OperatingForce, &r.ActivationTravel, &r.TotalTravel,
			&r.SoundProfile, &r.TriggerMethod, &r.Profile)

		result = append(result, r)
	}

	r.logger.LogTrace(fmt.Sprintf("result is %v", result))
	return result, nil
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
