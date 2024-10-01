package repo

import (
	"context"
	"kbswitch/internal/core/switches"
	"kbswitch/internal/core/switches/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(p *pgxpool.Pool) switches.Repo {
	return repo{
		pool: p,
	}
}

type repo struct {
	pool *pgxpool.Pool
}

// AddNew implements switches.Repo.
func (r repo) AddNew(context.Context, models.SwitchEntity) (*int, error) {
	panic("unimplemented")
}

// GetAll implements switches.Repo.
func (r repo) GetAll(ctx context.Context) ([]models.SwitchEntity, error) {
	result := make([]models.SwitchEntity, 0)
	query := `SELECT * FROM public.switches`

	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		// logger.Error(err.Error())
		return result, err
	}

	for rows.Next() {
		var r models.SwitchEntity
		err := rows.Scan(&r.ID, &r.Lifespan, &r.OperatingForce, &r.ActivationTravel,
			&r.TotalTravel, &r.Image, &r.Manufacturer, &r.Model, &r.ActuationType,
			&r.SoundProfile, &r.TriggerMethod, &r.Profile)
		if err != nil {
			// logger.Error(err.Error())
			return result, err
		}
		result = append(result, r)
	}

	if err = rows.Err(); err != nil {
		// logger.Error(err.Error())
		return result, err
	}

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
