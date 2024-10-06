package repo_test

import (
	"context"
	"encoding/json"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches/repo"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
)

func assertLogsEqual(method string, t *testing.T, want []string, got []string) {
	if !reflect.DeepEqual(want, got) && len(want) != len(got) {
		t.Errorf("in method %s: log check failed\nexpected %+v\ngot %v", method, want, got)
	}
}

func assertErrorReturned(method string, t *testing.T, want error, got error) {
	if want != nil && got == nil {
		t.Errorf("in method %s: expected error is not nil %v, when result returned nil: %v", method, want, got)
	}
}

func assertResultsEqual(method string, t *testing.T, want any, got any) {
	if !reflect.DeepEqual(want, got) {
		w, _ := json.Marshal(want)
		g, _ := json.Marshal(got)
		t.Errorf("in method %s: result check failed\nexpected %s\ngot %s", method, w, g)
	}
}

const (
	LogLvlInfo  = "I"
	LogLvlError = "E"
	LogLvlTrace = "T"
)

type fakeLogger struct {
	logs []string
}

// LogError implements logging.Logger.
func (f *fakeLogger) LogError(msg string) {
	f.logs = append(f.logs, LogLvlError)
}

// LogInfo implements logging.Logger.
func (f *fakeLogger) LogInfo(msg string) {
	f.logs = append(f.logs, LogLvlInfo)
}

// LogTrace implements logging.Logger.
func (f *fakeLogger) LogTrace(msg string) {
	f.logs = append(f.logs, LogLvlTrace)
}

type fakePool struct {
	getAllReturner func() (pgx.Rows, error)
	// getAllReturner func() (pgxmock.Rows, error)
}

// Exec implements database.DBPool.
func (f fakePool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	panic("unimplemented")
}

// Query implements database.DBPool.
func (f fakePool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return f.getAllReturner()
}

// QueryRow implements database.DBPool.
func (f fakePool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	panic("unimplemented")
}

func TestGetAll(t *testing.T) {
	cases := []struct {
		pool     fakePool
		logger   fakeLogger
		expected struct {
			res  []models.SwitchEntity
			err  error
			logs []string
		}
	}{
		{
			pool: fakePool{
				getAllReturner: func() (pgx.Rows, error) {
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
			logger: fakeLogger{},
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
				logs: []string{LogLvlTrace},
			},
		},
	}

	for _, tc := range cases {
		sut := repo.New(&tc.logger, tc.pool)
		got, _ := sut.GetAll(context.Background())
		want := tc.expected

		if !reflect.DeepEqual(want.res, got) {
			t.Errorf("want: %+v\ngot: %+v", want.res, got)
		}

		assertLogsEqual("GetAll", t, tc.expected.logs, tc.logger.logs)

	}
}
