package switches_test

import (
	"fmt"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/pkg/switches"
	"reflect"
	"testing"
)

type fakeRepo struct {
	getAllReturner func() ([]models.SwitchEntity, error)
}

// AddNew implements repositories.SwitchesRepo.
func (f fakeRepo) AddNew(models.SwitchEntity) error {
	panic("unimplemented")
}

// GetAll implements repositories.SwitchesRepo.
func (f fakeRepo) GetAll() ([]models.SwitchEntity, error) {
	return f.getAllReturner()
}

// GetByID implements repositories.SwitchesRepo.
func (f fakeRepo) GetByID(int) (*models.SwitchEntity, error) {
	panic("unimplemented")
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
