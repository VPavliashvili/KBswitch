package switches_test

import (
	"fmt"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches"
	"reflect"
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
func (f fakeRepo) Update(id int, req models.SwitchEntity) (*models.SwitchEntity, error) {
	return f.updateAction(id, req)
}

// Remove implements repositories.SwitchesRepo.
func (f fakeRepo) Remove(id int) error {
	return f.removeAction(id)
}

// AddNew implements repositories.SwitchesRepo.
func (f fakeRepo) AddNew(rb models.SwitchEntity) (*int, error) {
	return f.addNewAction(rb)
}

// GetAll implements repositories.SwitchesRepo.
func (f fakeRepo) GetAll() ([]models.SwitchEntity, error) {
	return f.getAllReturner()
}

// GetSingle implements repositories.SwitchesRepo.
func (f fakeRepo) GetSingle(id int) (*models.SwitchEntity, error) {
	return f.getSingleReturner(id)
}

// GetSingle implements repositories.SwitchesRepo.
func (f fakeRepo) GetID(brand, name string) (*int, error) {
	return f.getID(brand, name)
}

func TestRemove(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		brand    string
		name     string
		expected error
	}{
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
			},
			brand:    "test",
			name:     "test",
			expected: fmt.Errorf("resource with given brand and name not found"),
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, fmt.Errorf("test")
				},
			},
			brand:    "test",
			name:     "test",
			expected: fmt.Errorf("test"),
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), fmt.Errorf("test")
				},
			},
			brand:    "test",
			name:     "test",
			expected: fmt.Errorf("test"),
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
				removeAction: func(i int) error {
					return fmt.Errorf("test")
				},
			},
			brand:    "test",
			name:     "test",
			expected: fmt.Errorf("test"),
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
		unit := switches.GetService(tc.repo)
		err := unit.Remove(tc.brand, tc.name)

		if !reflect.DeepEqual(tc.expected, err) {
			t.Errorf("Remove error check failed\nexpected %v\ngot %v", tc.expected, err)
		}
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
			err error
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("resource with given brand and name not found"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return nil, fmt.Errorf("test")
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), fmt.Errorf("test")
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
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
				err error
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
					return nil, fmt.Errorf("test")
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Name: "tst"}, fmt.Errorf("test")
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(string, string) (*int, error) {
					return intptr(123), nil
				},
				updateAction: func(i int, se models.SwitchEntity) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Name: "tst"}, nil
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
				err error
			}{
				res: &models.Switch{Name: "tst"},
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.GetService(tc.repo)
		res, err := unit.Update(tc.in.brand, tc.in.name, tc.in.body)

		if !reflect.DeepEqual(tc.expected.err, err) {
			t.Errorf("Update error check failed\nexpected %v\ngot %v", tc.expected.err, err)
		}
		if !reflect.DeepEqual(tc.expected.res, res) {
			t.Errorf("Update result check failed\nexpected %v\ngot %v", tc.expected.res, res)
		}
	}
}

func TestAddNew(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		reqbody  models.SwitchRequestBody
		expected struct {
			res *int
			err error
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("switch with brand 'testb' and name 'testn' already exists"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, fmt.Errorf("test")
				},
			},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return intptr(123), fmt.Errorf("test")
				},
			},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
			},
		},
		{
			repo: fakeRepo{
				getID: func(s1, s2 string) (*int, error) {
					return nil, nil
				},
				addNewAction: func(se models.SwitchEntity) (*int, error) {
					return nil, fmt.Errorf("test")
				},
			},
			reqbody: models.SwitchRequestBody{Name: "testn", Brand: "testb"},
			expected: struct {
				res *int
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
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
				err error
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
				err error
			}{
				res: intptr(123),
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.GetService(tc.repo)
		res, err := unit.AddNew(tc.reqbody)

		if !reflect.DeepEqual(tc.expected.err, err) {
			t.Errorf("AddNew error check failed\nexpected %v\ngot %v", tc.expected.err, err)
		}
		if !reflect.DeepEqual(tc.expected.res, res) {
			t.Errorf("AddNew result check failed\nexpected %v\ngot %v", tc.expected.res, res)
		}
	}
}

func TestGetSingle(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		brand    string
		name     string
		expected struct {
			res *models.Switch
			err error
		}
	}{
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return nil, fmt.Errorf("test")
				},
				getID: func(s1, s2 string) (*int, error) { return intptr(123), nil }},
			expected: struct {
				res *models.Switch
				err error
			}{
				res: nil,
				err: fmt.Errorf("test"),
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
				err error
			}{
				res: nil,
				err: nil,
			},
		},
		{
			repo: fakeRepo{getID: func(s1, s2 string) (*int, error) {
				return nil, fmt.Errorf("test error from exists()")
			}},
			brand: "",
			name:  "",
			expected: struct {
				res *models.Switch
				err error
			}{
				res: nil,
				err: fmt.Errorf("test error from exists()"),
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("given combination of brand and name not found"),
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
				err error
			}{
				res: nil,
				err: fmt.Errorf("given combination of brand and name not found"),
			},
		},
		{
			repo: fakeRepo{
				getSingleReturner: func(int) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Name: "name", Brand: "brand"}, nil
				},
				getID: func(s1, s2 string) (*int, error) {
					return intptr(123), nil
				},
			},
			brand: "brand",
			name:  "name",
			expected: struct {
				res *models.Switch
				err error
			}{
				res: &models.Switch{Name: "name", Brand: "brand"},
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.GetService(tc.repo)
		res, err := unit.GetSingle(tc.brand, tc.name)

		if !reflect.DeepEqual(tc.expected.err, err) {
			t.Errorf("GetSingle error check failed\nexpected %v\ngot %v", tc.expected.err, err)
		}
		if !reflect.DeepEqual(tc.expected.res, res) {
			t.Errorf("GetSingle result check failed\nexpected %v\ngot %v", tc.expected.res, res)
		}
	}
}

func TestGetAll(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		expected struct {
			res []models.Switch
			err error
		}
	}{
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return nil, fmt.Errorf("test")
				},
			},
			expected: struct {
				res []models.Switch
				err error
			}{
				res: []models.Switch{},
				err: fmt.Errorf("test"),
			},
		},
		{
			repo: fakeRepo{
				getAllReturner: func() ([]models.SwitchEntity, error) {
					return []models.SwitchEntity{{Name: "testname", Brand: "idkbrand"}}, nil
				},
			},
			expected: struct {
				res []models.Switch
				err error
			}{
				res: []models.Switch{
					models.Switch{Name: "testname", Brand: "idkbrand"},
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
				err error
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
				err error
			}{
				res: []models.Switch{},
				err: nil,
			},
		},
	}

	for _, tc := range tcases {
		unit := switches.GetService(tc.repo)
		res, err := unit.GetAll()

		if !reflect.DeepEqual(tc.expected.err, err) {
			t.Errorf("GetAll error check failed\nexpected %v\ngot %v", tc.expected.err, err)
		}
		if !reflect.DeepEqual(tc.expected.res, res) {
			t.Errorf("GetAll result check failed\nexpected %v\ngot %v", tc.expected.res, res)
		}
	}
}
