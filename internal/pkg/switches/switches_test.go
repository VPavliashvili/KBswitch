package switches_test

import (
	"fmt"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches"
	"reflect"
	"testing"
)

type fakeRepo struct {
	getAllReturner    func() ([]models.SwitchEntity, error)
	getSingleReturner func(string, string) (*models.SwitchEntity, error)
	checkExists       func(string, string) (bool, error)
}

// AddNew implements repositories.SwitchesRepo.
func (f fakeRepo) AddNew(models.SwitchEntity) error {
	panic("unimplemented")
}

// GetAll implements repositories.SwitchesRepo.
func (f fakeRepo) GetAll() ([]models.SwitchEntity, error) {
	return f.getAllReturner()
}

// GetSingle implements repositories.SwitchesRepo.
func (f fakeRepo) GetSingle(brand, name string) (*models.SwitchEntity, error) {
	return f.getSingleReturner(brand, name)
}

func (f fakeRepo) Exists(brand, name string) (bool, error) {
	return f.checkExists(brand, name)
}

func TestGetByID(t *testing.T) {
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
				getSingleReturner: func(s1, s2 string) (*models.SwitchEntity, error) {
					return nil, fmt.Errorf("test")
				},
				checkExists: func(s1, s2 string) (bool, error) { return true, nil }},
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
				getSingleReturner: func(brand, name string) (*models.SwitchEntity, error) {
					return nil, nil
				},
				checkExists: func(s1, s2 string) (bool, error) { return true, nil },
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
			repo: fakeRepo{checkExists: func(s1, s2 string) (bool, error) {
				return false, fmt.Errorf("test error from exists()")
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
				getSingleReturner: func(string, string) (*models.SwitchEntity, error) {
					return nil, fmt.Errorf("given combination of brand and name not found")
				},
				checkExists: func(s1, s2 string) (bool, error) {
					return false, nil
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
				getSingleReturner: func(string, string) (*models.SwitchEntity, error) {
					return nil, nil
				},
				checkExists: func(s1, s2 string) (bool, error) {
					return false, nil
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
				getSingleReturner: func(s1, s2 string) (*models.SwitchEntity, error) {
					return &models.SwitchEntity{Name: "name", Brand: "brand"}, nil
				},
				checkExists: func(s1, s2 string) (bool, error) {
					return true, nil
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
		unit := switches.New(tc.repo)
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
		unit := switches.New(tc.repo)
		res, err := unit.GetAll()

		if !reflect.DeepEqual(tc.expected.err, err) {
			t.Errorf("GetAll error check failed\nexpected %v\ngot %v", tc.expected.err, err)
		}
		if !reflect.DeepEqual(tc.expected.res, res) {
			t.Errorf("GetAll result check failed\nexpected %v\ngot %v", tc.expected.res, res)
		}
	}
}
