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

type fakeLogger struct{}

// LogError implements logging.Logger.
func (f fakeLogger) LogError(msg string) {
	panic("unimplemented")
}

// LogInfo implements logging.Logger.
func (f fakeLogger) LogInfo(msg string) {
	panic("unimplemented")
}

// LogTrace implements logging.Logger.
func (f fakeLogger) LogTrace(msg string) {
	panic("unimplemented")
}

func TestRemove(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		brand    string
		name     string
		expected *common.AppError
	}{
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
			},
			brand:    "test",
			name:     "test",
			expected: &switches.ErrNoSwitch,
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, errTest
				},
			},
			brand:    "test",
			name:     "test",
			expected: common.Wrap(errTest),
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), errTest
				},
			},
			brand:    "test",
			name:     "test",
			expected: common.Wrap(errTest),
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
			brand:    "test",
			name:     "test",
			expected: common.Wrap(errTest),
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
			brand:    "test",
			name:     "test",
			expected: nil,
		},
	}

	for _, tc := range tcases {
		unit := switches.New(nil, tc.repo)
		err := unit.Remove(context.Background(), tc.brand, tc.name)

		assertErrorsEqual("Remove", t, tc.expected, err)
	}
}

func TestUpdate(t *testing.T) {
	tcases := []struct {
		repo fakeRepo
		in   struct {
			brand string
			name  string
			body  models.SwitchRequestBody
		}
		expected struct {
			res *models.Switch
			err *common.AppError
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
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: &switches.ErrNoSwitch,
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
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: nil,
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
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
				res *models.Switch
				err *common.AppError
			}{
				res: &models.Switch{Name: "tst"},
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(nil, tc.repo)
		res, err := unit.Update(context.Background(), tc.in.brand, tc.in.name, tc.in.body)

		assertErrorsEqual("Update", t, tc.expected.err, err)
		assertResultsEqual("Update", t, tc.expected.res, res)
	}
}

func TestAddNew(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		reqbody  models.SwitchRequestBody
		expected struct {
			res *int
			err *common.AppError
		}
	}{
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
			},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err *common.AppError
			}{
				res: nil,
				err: &switches.ErrAlreadyExists,
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, errTest
				},
			},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err *common.AppError
			}{
				res: nil,
				err: nil,
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
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err *common.AppError
			}{
				res: intptr(123),
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(nil, tc.repo)
		res, err := unit.AddNew(context.Background(), tc.reqbody)

		assertErrorsEqual("AddNew", t, tc.expected.err, err)
		assertResultsEqual("AddNew", t, tc.expected.res, res)
	}
}

func TestGetSingle(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		brand    string
		name     string
		expected struct {
			res *models.Switch
			err *common.AppError
		}
	}{
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, errTest
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil }},
			expected: struct {
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, nil
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil },
			},

			expected: struct {
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: &switches.ErrErrorMissing,
			},
		},
		{
			repo: fakeRepo{getID: func(s1, s2 string) (*int, error) {
				return nil, errTest
			}},
			brand: "",
			name:  "",
			expected: struct {
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: common.Wrap(errTest),
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
			brand: "bad brand",
			name:  "or bad name",
			expected: struct {
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: &switches.ErrNoSwitch,
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
			brand: "bad brand",
			name:  "or bad name",
			expected: struct {
				res *models.Switch
				err *common.AppError
			}{
				res: nil,
				err: &switches.ErrNoSwitch,
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
			brand: "brand",
			name:  "name",
			expected: struct {
				res *models.Switch
				err *common.AppError
			}{
				res: &models.Switch{Name: "name", Brand: "brand"},
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(nil, tc.repo)
		res, err := unit.GetSingle(context.Background(), tc.brand, tc.name)

		assertErrorsEqual("GetSingle", t, tc.expected.err, err)
		assertResultsEqual("GetSingle", t, tc.expected.res, res)
	}
}

func TestGetAll(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		logger   fakeLogger
		expected struct {
			res []models.Switch
			err *common.AppError
		}
	}{
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, errTest
				},
			},
			expected: struct {
				res []models.Switch
				err *common.AppError
			}{
				res: []models.Switch{},
				err: common.Wrap(errTest),
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{{Model: "testname", Manufacturer: "idkbrand"}}, nil
				},
			},
			expected: struct {
				res []models.Switch
				err *common.AppError
			}{
				res: []models.Switch{
					{Name: "testname", Brand: "idkbrand"},
				},
				err: nil,
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, nil
				},
			},
			expected: struct {
				res []models.Switch
				err *common.AppError
			}{
				res: []models.Switch{},
				err: nil,
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{}, nil
				},
			},
			expected: struct {
				res []models.Switch
				err *common.AppError
			}{
				res: []models.Switch{},
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.New(tc.logger, tc.repo)
		res, err := unit.GetAll(context.Background())

		assertErrorsEqual("GetAll", t, tc.expected.err, err)
		assertResultsEqual("GetAll", t, tc.expected.res, res)
	}
}
