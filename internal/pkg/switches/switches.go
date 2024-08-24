package switches

import (
	"fmt"
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/core/switches/repositories"
)

func New(repo repositories.SwitchesRepo) service {
	return service{repo: repo}
}

type service struct {
	repo repositories.SwitchesRepo
}

func (s service) AddNew(reqbody models.SwitchRequestBody) (*int, error) {
	exists, err := s.repo.Exists(reqbody.Brand, reqbody.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("switch with brand '%s' and name '%s' already exists", reqbody.Brand, reqbody.Name)
	}

	entity := models.SwitchEntity(reqbody)

	resp, err := s.repo.AddNew(entity)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s service) GetAll() ([]models.Switch, error) {
	res := []models.Switch{}
	resp, err := s.repo.GetAll()

	for _, item := range resp {
		// s := models.Switch{
		// 	Brand:            item.Brand,
		// 	ActuationType:    item.ActuationType,
		// 	Lifespan:         item.Lifespan,
		// 	Name:             item.Name,
		// 	Image:            item.Image,
		// 	OperatingForce:   item.OperatingForce,
		// 	ActivationTravel: item.ActivationTravel,
		// 	TotalTravel:      item.TotalTravel,
		// 	SoundProfile:     item.SoundProfile,
		// 	TriggerMethod:    item.TriggerMethod,
		// 	Profile:          item.Profile,
		// }
		s := models.Switch(item)

		res = append(res, s)
	}

	return res, err
}

func (s service) GetSingle(brand, name string) (*models.Switch, error) {
	exists, err := s.repo.Exists(brand, name)

	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("given combination of brand and name not found")
	}

	resp, err := s.repo.GetSingle(brand, name)
	if resp == nil {
		return nil, err
	}
	res := models.Switch(*resp)

	return &res, err
}
