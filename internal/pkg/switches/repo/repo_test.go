package repo_test

import (
	"context"
	"fmt"
	"kbswitch/internal/core/common/tests"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches/repo"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeRow struct {
	scan func(...any) error
}

func (f fakeRow) Scan(dest ...any) error {
	return f.scan(dest)
}

type fakeRows struct {
	next func() bool
	scan func(...any) error
}

var connclosed bool

// Close implements pgx.Rows.
func (f fakeRows) Close() {
	connclosed = true
}

// CommandTag implements pgx.Rows.
func (f fakeRows) CommandTag() pgconn.CommandTag {
	panic("unimplemented")
}

// Conn implements pgx.Rows.
func (f fakeRows) Conn() *pgx.Conn {
	panic("unimplemented")
}

// Err implements pgx.Rows.
func (f fakeRows) Err() error {
	panic("unimplemented")
}

// FieldDescriptions implements pgx.Rows.
func (f fakeRows) FieldDescriptions() []pgconn.FieldDescription {
	panic("unimplemented")
}

// Next implements pgx.Rows.
func (f fakeRows) Next() bool {
	return f.next()
}

// RawValues implements pgx.Rows.
func (f fakeRows) RawValues() [][]byte {
	panic("unimplemented")
}

// Scan implements pgx.Rows.
func (f fakeRows) Scan(dest ...any) error {
	return f.scan(dest)
}

// Values implements pgx.Rows.
func (f fakeRows) Values() ([]any, error) {
	panic("unimplemented")
}

type fakePool struct {
	queryFunc        func() (pgx.Rows, error)
	querySingleFunc  func(int) pgx.Row
	getSingleIdParam int
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
	return f.querySingleFunc(f.getSingleIdParam)
}

func assign[T int | float64 | string | []byte](source any, value T) {
	if x, ok := source.(*T); ok {
		*x = value
	} else {
		msg := fmt.Sprintf("scan assignment not ok. value '%v' of type %T can't be cast on source with type %s",
			value, value, reflect.TypeOf(source))
		panic(msg)
	}
}

func scanObject(id, ls, of int, att, tt float64, img []byte, mn, mm, at, sp, tm, p string, dest ...any) {
	target := dest[0].([]any)

	assign(target[0], id)
	assign(target[1], ls)
	assign(target[2], of)
	assign(target[3], att)
	assign(target[4], tt)
	assign(target[5], img)
	assign(target[6], mn)
	assign(target[7], mm)
	assign(target[8], at)
	assign(target[9], sp)
	assign(target[10], tm)
	assign(target[11], p)

	dest[0] = target
}

func TestGetSingle(t *testing.T) {
	cases := []struct {
		pool     fakePool
		logger   tests.FakeLogger
		expected struct {
			res  *models.SwitchEntity
			err  error
			logs []string
		}
	}{
		{
			pool: fakePool{
				getSingleIdParam: 123,
				querySingleFunc: func(i int) pgx.Row {
					row := fakeRow{
						scan: func(dest ...any) error {
							scanObject(1, 10, 30, 30, 30, []byte{1, 1}, "mn", "mm", "at", "sp", "tm", "p", dest...)
							return nil
						},
					}
					return row
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  *models.SwitchEntity
				err  error
				logs []string
			}{
				res: &models.SwitchEntity{
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
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			pool: fakePool{
				getSingleIdParam: 123,
				querySingleFunc: func(i int) pgx.Row {
					row := fakeRow{
						scan: func(a ...any) error {
							return pgx.ErrNoRows
						},
					}
					return row
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  *models.SwitchEntity
				err  error
				logs []string
			}{
				res:  nil,
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			pool: fakePool{
				getSingleIdParam: 123,
				querySingleFunc: func(i int) pgx.Row {
					row := fakeRow{
						scan: func(a ...any) error {
							return tests.ErrTest
						},
					}
					return row
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  *models.SwitchEntity
				err  error
				logs []string
			}{
				res:  nil,
				err:  tests.ErrTest,
				logs: []string{tests.LogLvlError},
			},
		},
	}

	for _, tc := range cases {
		sut := repo.New(&tc.logger, tc.pool)
		got, err := sut.GetSingle(context.Background(), tc.pool.getSingleIdParam)

		tests.AssertHasError("GetSingle", t, tc.expected.err, err)
		tests.AssertResultsEqual("GetSingle", t, tc.expected.res, got)
		tests.AssertLogsEqual("GetSingle", t, tc.expected.logs, tc.logger.Logs)
	}
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
					counter := 0
					elements := 2
					rows := fakeRows{
						next: func() bool {
							res := elements > counter
							counter++
							return res
						},
						scan: func(dest ...any) error {
							switch counter {
							case 1:
								scanObject(1, 10, 30, 30, 30, []byte{1, 1}, "mn", "mm", "at", "sp", "tm", "p", dest...)
							case 2:
								scanObject(2, 20, 40, 40, 40, []byte{2, 2}, "mn2", "mm2", "at2", "sp2", "tm2", "p2", dest...)
							}
							return nil
						},
					}

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
					rows := fakeRows{
						next: func() bool {
							// i.e no elements returned from db
							return false
						},
						scan: func(dest ...any) error {
							// this won't even run because scan only could run after next
							return nil
						},
					}

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
		{
			pool: fakePool{
				queryFunc: func() (pgx.Rows, error) {
					counter := 0
					rowsCount := 1
					rows := fakeRows{
						next: func() bool {
							res := rowsCount > counter
							counter++
							return res
						},
						scan: func(...any) error {
							return tests.ErrTest
						},
					}
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
				err:  tests.ErrTest,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			pool: fakePool{
				queryFunc: func() (pgx.Rows, error) {
					counter := 0
					rowsCount := 1
					rows := fakeRows{
						next: func() bool {
							res := rowsCount > counter
							counter++
							return res
						},
						scan: func(...any) error {
							return nil
						},
					}
					defer func() {
						if !connclosed {
							panic("GETALL DID NOT CLOSE THE CONNECTION")
						}
						// reset for potentially other test cases
						connclosed = false
					}()
					return rows, nil
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.SwitchEntity
				err  error
				logs []string
			}{
				res:  []models.SwitchEntity{{}},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
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
