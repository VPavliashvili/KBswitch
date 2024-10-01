package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Key string

// list of pool map keys
var PoolKey = struct {
	Switches Key
}{
	Switches: "switches",
}

var pool = make(map[Key]*pgxpool.Pool)

// appends a new pool with a given key
func Append(k Key, p *pgxpool.Pool) {
	_, exists := pool[k]
	if !exists {
		pool[k] = p
	}
}

// gets pool by key
func Get(k Key) *pgxpool.Pool {
	p, exists := pool[k]
	if exists {
		return p
	}
	panic(fmt.Sprintf("pool does not contain object with key: %s", k))
}

// checks database availability
func Ping(ctx context.Context, k Key) error {
	p := Get(k)
	err := p.Ping(ctx)
	return err
}
