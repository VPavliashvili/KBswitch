package switches_test

import (
	"context"
	"encoding/json"
	"fmt"
	"kbswitch/internal/core/common"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches"
	"reflect"
	"testing"
)

func assertLogsEqual(method string, t *testing.T, want []string, got []string) {
	if !reflect.DeepEqual(want, got) && len(want) != len(got) {
		t.Errorf("in method %s: log check failed\nexpected %+v\ngot %v", method, want, got)
	}
}

func assertErrorsEqual(method string, t *testing.T, want *common.AppError, got *common.AppError) {
	if want == nil && got != nil {
		t.Errorf("in method %s: expected error equals to nil, when error returned: %v", method, got)
	} else if want != nil {
		et := want.Errtype.Error() == got.Errtype.Error()
		er := want.Reason.Error() == got.Reason.Error()
		equals := et && er

		if !equals {
			t.Errorf("in method %s: error check failed\nexpected type: %v\ngot type: %v\nexpected reason: %v\ngot reason: %v",
				method, want.Errtype, got.Errtype, want.Reason, got.Reason)
		}
	}
}

func assertResultsEqual(method string, t *testing.T, want any, got any) {
	if !reflect.DeepEqual(want, got) {
		w, _ := json.Marshal(want)
		g, _ := json.Marshal(got)
		t.Errorf("in method %s: result check failed\nexpected %s\ngot %s", method, w, g)
	}
}

var errTest = fmt.Errorf("test")

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

func TestRemove(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   fakeLogger
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
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, errTest
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), errTest
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
				removeAction: func(i int) error {
					return errTest
				},
			},
			brand: "test",
			name:  "test",
			expected: struct {
				err  *common.AppError
				logs []string
			}{
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
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
				logs: []string{LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		err := unit.Remove(context.Background(), tc.brand, tc.name)

		assertErrorsEqual("Remove", t, tc.expected.err, err)
		assertLogsEqual("Remove", t, tc.expected.logs, tc.logger.logs)
	}
}

func TestUpdate(t *testing.T) {
	tcases := []struct {
		repo   fakeRepo
		logger fakeLogger
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
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return nil, errTest
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
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), errTest
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
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
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
				logs: []string{},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return nil, errTest
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
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Model: "tst"}, errTest
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
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
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
				logs: []string{LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.Update(context.Background(), tc.in.brand, tc.in.name, tc.in.body)

		assertErrorsEqual("Update", t, tc.expected.err, err)
		assertResultsEqual("Update", t, tc.expected.res, res)
		assertLogsEqual("Update", t, tc.expected.logs, tc.logger.logs)
	}
}

func TestAddNew(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		reqbody  models.SwitchRequestBody
		logger   fakeLogger
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
			logger:  fakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrAlreadyExists,
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, errTest
				},
			},
			logger:  fakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return intptr(123), errTest
				},
			},
			logger:  fakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return nil, errTest
				},
			},
			logger:  fakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
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
			logger:  fakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  nil,
				logs: []string{LogLvlTrace},
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return intptr(123), nil
				},
			},
			logger:  fakeLogger{},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res  *int
				err  *common.AppError
				logs []string
			}{
				res:  intptr(123),
				err:  nil,
				logs: []string{LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.AddNew(context.Background(), tc.reqbody)

		assertErrorsEqual("AddNew", t, tc.expected.err, err)
		assertResultsEqual("AddNew", t, tc.expected.res, res)
		assertLogsEqual("AddNew", t, tc.expected.logs, tc.logger.logs)
	}
}

func TestGetSingle(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   fakeLogger
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
					return nil, errTest
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil }},
			logger: fakeLogger{},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, nil
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil },
			},
			logger: fakeLogger{},
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrErrorMissing,
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{getID: func(s1, s2 string) (*int, error) {
				return nil, errTest
			}},
			logger: fakeLogger{},
			brand:  "",
			name:   "",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
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
			logger: fakeLogger{},
			brand:  "bad brand",
			name:   "or bad name",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrNoSwitch,
				logs: []string{LogLvlError},
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
			logger: fakeLogger{},
			brand:  "bad brand",
			name:   "or bad name",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  nil,
				err:  &switches.ErrNoSwitch,
				logs: []string{LogLvlError},
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
			logger: fakeLogger{},
			brand:  "brand",
			name:   "name",
			expected: struct {
				res  *models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  &models.Switch{Name: "name", Brand: "brand"},
				err:  nil,
				logs: []string{LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.GetSingle(context.Background(), tc.brand, tc.name)

		assertErrorsEqual("GetSingle", t, tc.expected.err, err)
		assertResultsEqual("GetSingle", t, tc.expected.res, res)
		assertLogsEqual("GetSingle", t, tc.expected.logs, tc.logger.logs)
	}
}

func TestGetAll(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   fakeLogger
		expected struct {
			res  []models.Switch
			err  *common.AppError
			logs []string
		}
	}{
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, errTest
				},
			},
			logger: fakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  []models.Switch{},
				err:  common.Wrap(errTest),
				logs: []string{LogLvlError},
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{{Model: "testname", Manufacturer: "idkbrand"}}, nil
				},
			},
			logger: fakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res: []models.Switch{
					{Name: "testname", Brand: "idkbrand"},
				},
				err:  nil,
				logs: []string{LogLvlTrace},
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, nil
				},
			},
			logger: fakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  []models.Switch{},
				err:  nil,
				logs: []string{LogLvlTrace},
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{}, nil
				},
			},
			logger: fakeLogger{},
			expected: struct {
				res  []models.Switch
				err  *common.AppError
				logs []string
			}{
				res:  []models.Switch{},
				err:  nil,
				logs: []string{LogLvlTrace},
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(&tc.logger, tc.repo)
		res, err := unit.GetAll(context.Background())

		assertErrorsEqual("GetAll", t, tc.expected.err, err)
		assertResultsEqual("GetAll", t, tc.expected.res, res)
		assertLogsEqual("GetAll", t, tc.expected.logs, tc.logger.logs)
	}
}
