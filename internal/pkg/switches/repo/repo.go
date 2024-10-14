package repo

import (
	"context"
	"fmt"
	"kbswitch/internal/app"
	"kbswitch/internal/core/common/database"
	"kbswitch/internal/core/common/logging"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"

	"github.com/jackc/pgx/v5"
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
		var s models.SwitchEntity
		err = rows.Scan(&s.ID, &s.Lifespan, &s.OperatingForce, &s.ActivationTravel,
			&s.TotalTravel, &s.Image, &s.Manufacturer, &s.Model, &s.ActuationType,
			&s.SoundProfile, &s.TriggerMethod, &s.Profile)
		if err != nil {
			r.logger.LogError(err.Error())
			return []models.SwitchEntity{}, err
		}

		result = append(result, s)
	}

	r.logger.LogTrace(fmt.Sprintf("result is %v", result))
	return result, nil
}

// GetID implements switches.Repo.
func (r repo) GetID(ctx context.Context, brand string, name string) (*int, error) {
	panic("unimplemented")
}

// GetSingle implements switches.Repo.
func (r repo) GetSingle(ctx context.Context, id int) (*models.SwitchEntity, error) {
	query := `SELECT * FROM public.switches WHERE id=$1`
	row := r.pool.QueryRow(ctx, query)

	var s models.SwitchEntity
	err := row.Scan(&s.ID, &s.Lifespan, &s.OperatingForce, &s.ActivationTravel,
		&s.TotalTravel, &s.Image, &s.Manufacturer, &s.Model, &s.ActuationType,
		&s.SoundProfile, &s.TriggerMethod, &s.Profile)
	if err == pgx.ErrNoRows {
		r.logger.LogTrace(fmt.Sprintf("no result found for id: %v", id))
		return nil, nil
	}
	if err != nil {
		r.logger.LogError(fmt.Sprintf("query error: %s", err.Error()))
		return nil, err
	}

	r.logger.LogTrace(fmt.Sprintf("result is %v", s))
	return &s, nil
}

// Remove implements switches.Repo.
func (r repo) Remove(context.Context, int) error {
	panic("unimplemented")
}

// Update implements switches.Repo.
func (r repo) Update(context.Context, int, models.SwitchEntity) (*models.SwitchEntity, error) {
	panic("unimplemented")
}
