package repo_test

import (
	"context"
	"kbswitch/internal/core/common/tests"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches/repo"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
)

type fakePool struct {
	queryFunc func() (pgx.Rows, error)
}

// Exec implements database.DBPool.
func (f fakePool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	panic("unimplemented")
}

// Query implements database.DBPool.
func (f fakePool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return f.queryFunc()
}

// QueryRow implements database.DBPool.
func (f fakePool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	panic("unimplemented")
}

func TestGetAll(t *testing.T) {
	cases := []struct {
		pool     fakePool
		logger   tests.FakeLogger
		expected struct {
			res  []models.SwitchEntity
			err  error
			logs []string
		}
	}{
		{
			pool: fakePool{
				queryFunc: func() (pgx.Rows, error) {
					c, _ := pgxmock.NewConn()
					defer c.Close(context.Background())

					columns := []string{
						"id", "manufacturer", "actuationType",
						"lifespan", "model", "image", "operatingForce",
						"activationTravel", "totalTravel", "soundProfile",
						"triggerMethod", "profile",
					}

					rows := c.NewRows(columns).AddRow(1, "mn", "at", 10, "mm", []byte{1, 1}, 30, float64(30), float64(30), "sp", "tm", "p").
						AddRow(2, "mn2", "at2", 20, "mm2", []byte{2, 2}, 40, float64(40), float64(40), "sp2", "tm2", "p2").
						Kind()

					return rows, nil
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.SwitchEntity
				err  error
				logs []string
			}{
				res: []models.SwitchEntity{
					{
						ID:               1,
						Manufacturer:     "mn",
						ActuationType:    "at",
						Lifespan:         10,
						Model:            "mm",
						Image:            []byte{1, 1},
						OperatingForce:   30,
						ActivationTravel: 30,
						TotalTravel:      30,
						SoundProfile:     "sp",
						TriggerMethod:    "tm",
						Profile:          "p",
					},
					{
						ID:               2,
						Manufacturer:     "mn2",
						ActuationType:    "at2",
						Lifespan:         20,
						Model:            "mm2",
						Image:            []byte{2, 2},
						OperatingForce:   40,
						ActivationTravel: 40,
						TotalTravel:      40,
						SoundProfile:     "sp2",
						TriggerMethod:    "tm2",
						Profile:          "p2",
					},
				},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			pool: fakePool{
				queryFunc: func() (pgx.Rows, error) {
					c, _ := pgxmock.NewConn()
					defer c.Close(context.Background())

					columns := []string{
						"id", "manufacturer", "actuationType",
						"lifespan", "model", "image", "operatingForce",
						"activationTravel", "totalTravel", "soundProfile",
						"triggerMethod", "profile",
					}
					rows := c.NewRows(columns).Kind()

					return rows, nil
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.SwitchEntity
				err  error
				logs []string
			}{
				res:  []models.SwitchEntity{},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			pool: fakePool{
				queryFunc: func() (pgx.Rows, error) {
					return nil, tests.ErrTest
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.SwitchEntity
				err  error
				logs []string
			}{
				res:  []models.SwitchEntity{},
				err:  tests.ErrTest,
				logs: []string{tests.LogLvlError},
			},
		},
	}

	for _, tc := range cases {
		sut := repo.New(&tc.logger, tc.pool)
		got, err := sut.GetAll(context.Background())

		tests.AssertHasError("GetAll", t, tc.expected.err, err)
		tests.AssertResultsEqual("GetAll", t, tc.expected.res, got)
		tests.AssertLogsEqual("GetAll", t, tc.expected.logs, tc.logger.Logs)
	}
}
