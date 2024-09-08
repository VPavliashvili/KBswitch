package middlewares

import (
	"fmt"
	"kbswitch/internal/app"
	"kbswitch/internal/core/common/database"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPgxPool(key database.Key, cfg app.DbConfig) func(http http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Db)

			var err error
			pool, err := pgxpool.New(r.Context(), dbUrl)
			if err != nil {
				// logger.Fatal(fmt.Sprintf("FATAL when creating pgx pool -> %s", err.Error()))
				err = fmt.Errorf("could not create pgx pool\n" + err.Error())
				panic(err.Error())
			}

			database.Append(key, pool)
			next.ServeHTTP(w, r)
		}

		println("pgxPool initialized, key: " + key)
		return http.HandlerFunc(fn)
	}
}
