package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool = make(map[string]*pgxpool.Pool)

func Append(key string, p *pgxpool.Pool) {
	_, exists := pool[key]
	if !exists {
		pool[key] = p
	}
}

func Get(key string) *pgxpool.Pool {
	p, exists := pool[key]
	if exists {
		return p
	}
	panic(fmt.Sprintf("pool does not contain object with key: %s", key))
}

func Ping(ctx context.Context, key string) error {
	p := Get(key)
	err := p.Ping(ctx)
	return err
}
