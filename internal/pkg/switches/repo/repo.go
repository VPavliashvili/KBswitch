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
	// pool, err := database.NewPool(ctx, r.cfg)
	// if err != nil {
	// 	// reqId := ctx.Value(logger.LogIDKey)
	// 	// logger.Error(fmt.Sprintf("on requestId: %s, error happened: %s", reqId, err.Error()))
	// 	return nil, err
	// }
	//
	// result := make([]models.SwitchEntity, 0)
	// query := `SELECT * FROM public.switches`
	//
	// rows, err := pool.Query(context.Background(), query)
	// if err != nil {
	// 	return result, err
	// }
	//
	// for rows.Next() {
	// 	var r models.SwitchEntity
	// 	err := rows.Scan(&r.ID, &r.Lifespan, &r.OperatingForce, &r.ActivationTravel,
	// 		&r.TotalTravel, &r.Image, &r.Manufacturer, &r.Model, &r.ActuationType,
	// 		&r.SoundProfile, &r.TriggerMethod, &r.Profile)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result = append(result, r)
	// }
	//
	// if err = rows.Err(); err != nil {
	// 	return result, err
	// }
	//
	// return result, nil

	result := make([]models.SwitchEntity, 0)
	query := `SELECT * FROM public.switches`

	rows, _ := r.pool.Query(ctx, query)
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
