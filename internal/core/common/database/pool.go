package database

import (
	"context"
	"fmt"
	"kbswitch/internal/app"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg app.DbConfig) (*pgxpool.Pool, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Db)

	pool, connParsingErr := pgxpool.New(ctx, dbUrl)
	if connParsingErr != nil {
		connParsingErr = fmt.Errorf("could not create pgx pool %s", connParsingErr.Error())
		return nil, connParsingErr
	}

	return pool, nil
}
