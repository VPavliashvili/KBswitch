package switches_test

import (
	"context"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/common/tests"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches"
	"testing"
)

func intptr(x int) *int {
	return &x
}

type fakeRepo struct {
	getID             func(string, string) (*int, error)
	getAllReturner    func() ([]models.SwitchEntity, error)
	getSingleReturner func(int) (*models.SwitchEntity, error)
	addNewAction      func(models.SwitchEntity) (*int, error)
	removeAction      func(int) error
	updateAction      func(int, models.SwitchEntity) (*models.SwitchEntity, error)
}

// Update implements repositories.SwitchesRepo.
func (f fakeRepo) Update(ctx context.Context, id int, req models.SwitchEntity) (*models.SwitchEntity, error) {
	return f.updateAction(id, req)
}

// Remove implements repositories.SwitchesRepo.
func (f fakeRepo) Remove(ctx context.Context, id int) error {
	return f.removeAction(id)
}

// AddNew implements repositories.SwitchesRepo.
func (f fakeRepo) AddNew(ctx context.Context, rb models.SwitchEntity) (*int, error) {
	return f.addNewAction(rb)
}

// GetAll implements repositories.SwitchesRepo.
func (f fakeRepo) GetAll(ctx context.Context) ([]models.SwitchEntity, error) {
	return f.getAllReturner()
}

// GetSingle implements repositories.SwitchesRepo.
func (f fakeRepo) GetSingle(ctx context.Context, id int) (*models.SwitchEntity, error) {
	return f.getSingleReturner(id)
}

// GetSingle implements repositories.SwitchesRepo.
func (f fakeRepo) GetID(ctx context.Context, brand, name string) (*int, error) {
	return f.getID(brand, name)
}

func TestRemove(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   tests.FakeLogger
		brand    string
		name     string
		expected struct {
			err  *common.AppError
			logs []string
		}
	}{
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  &switches.ErrNoSwitch,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, tests.ErrTest
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), tests.ErrTest
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
				removeAction: func(i int) error {
					return tests.ErrTest
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
				removeAction: func(i int) error {
					return nil
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		err := unit.Remove(context.Background(), tc.brand, tc.name)

		tests.AssertErrorsEqual("Remove", t, tc.expected.err, err)
		tests.AssertLogsEqual("Remove", t, tc.expected.logs, tc.logger.Logs)
	}
}

func TestUpdate(t *testing.T) {
	tcases := []struct {
		repo   fakeRepo
		logger tests.FakeLogger
		in     struct {
			brand string
			name  string
			body  models.SwitchRequestBody
		}
		expected struct {
			res  *models.Switch
			err  *common.AppError
			logs []string
		}
	}{
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return nil, nil
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrNoSwitch,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return nil, tests.ErrTest
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), tests.ErrTest
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return nil, nil
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return nil, tests.ErrTest
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Model: "tst"}, tests.ErrTest
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Model: "tst"}, nil
				},
			},
			in: struct {
				brand string
				name  string
				body  models.SwitchRequestBody
			}{
				brand: "test",
				name:  "test",
				body:  models.SwitchRequestBody{Brand: "newb", Name: "newn"},
			},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  &models.Switch{Name: "tst"},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.Update(context.Background(), tc.in.brand, tc.in.name, tc.in.body)

		tests.AssertErrorsEqual("Update", t, tc.expected.err, err)
		tests.AssertResultsEqual("Update", t, tc.expected.res, res)
		tests.AssertLogsEqual("Update", t, tc.expected.logs, tc.logger.Logs)
	}
}

func TestAddNew(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		reqbody  models.SwitchRequestBody
		logger   tests.FakeLogger
		expected struct {
			res  *int
			err  *common.AppError
			logs []string
		}
	}{
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
			},
			logger:  tests.FakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrAlreadyExists,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, tests.ErrTest
				},
			},
			logger:  tests.FakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return intptr(123), tests.ErrTest
				},
			},
			logger:  tests.FakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return nil, tests.ErrTest
				},
			},
			logger:  tests.FakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return nil, nil
				},
			},
			logger:  tests.FakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{repo: fakeRepo{
			getID: func(s1, s2 string) (*int, error) {
				return nil, nil
			},
			addNewAction: func(se models.SwitchEntity) (*int, error) {
				return intptr(123), nil
			},
		}, logger: tests.FakeLogger{}, reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"}, expected: struct {
			res  *int
			err  *common.AppError
			logs []string
		}{
			res:  intptr(123),
			err:  nil,
			logs: []string{tests.LogLvlTrace},
		}},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.AddNew(context.Background(), tc.reqbody)

		tests.AssertErrorsEqual("AddNew", t, tc.expected.err, err)
		tests.AssertResultsEqual("AddNew", t, tc.expected.res, res)
		tests.AssertLogsEqual("AddNew", t, tc.expected.logs, tc.logger.Logs)
	}
}

func TestGetSingle(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   tests.FakeLogger
		brand    string
		name     string
		expected struct {
			res  *models.Switch
			err  *common.AppError
			logs []string
		}
	}{
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, tests.ErrTest
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil }},
			logger: tests.FakeLogger{},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, nil
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil },
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrErrorMissing,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{getID: func(s1, s2 string) (*int, error) {
				return nil, tests.ErrTest
			}},
			logger: tests.FakeLogger{},
			brand:  "",
			name:   "",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, nil
				},
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
			},
			logger: tests.FakeLogger{},
			brand:  "bad brand",
			name:   "or bad name",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrNoSwitch,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, nil
				},
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
			},
			logger: tests.FakeLogger{},
			brand:  "bad brand",
			name:   "or bad name",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrNoSwitch,
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Model: "name", Manufacturer: "brand"}, nil
				},
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
			},
			logger: tests.FakeLogger{},
			brand:  "brand",
			name:   "name",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  &models.Switch{Name: "name", Brand: "brand"},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.GetSingle(context.Background(), tc.brand, tc.name)

		tests.AssertErrorsEqual("GetSingle", t, tc.expected.err, err)
		tests.AssertResultsEqual("GetSingle", t, tc.expected.res, res)
		tests.AssertLogsEqual("GetSingle", t, tc.expected.logs, tc.logger.Logs)
	}
}

func TestGetAll(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   tests.FakeLogger
		expected struct {
			res  []models.Switch
			err  *common.AppError
			logs []string
		}
	}{
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, tests.ErrTest
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  []models.Switch{},
				err:  common.Wrap(tests.ErrTest),
				logs: []string{tests.LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{{Model: "testname", Manufacturer: "idkbrand"}}, nil
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res: []models.Switch{
					{Name: "testname", Brand: "idkbrand"},
				},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, nil
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  []models.Switch{},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{}, nil
				},
			},
			logger: tests.FakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  []models.Switch{},
				err:  nil,
				logs: []string{tests.LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.GetAll(context.Background())

		tests.AssertErrorsEqual("GetAll", t, tc.expected.err, err)
		tests.AssertResultsEqual("GetAll", t, tc.expected.res, res)
		tests.AssertLogsEqual("GetAll", t, tc.expected.logs, tc.logger.Logs)
	}
}
