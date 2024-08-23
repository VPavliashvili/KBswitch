package switches

import (
	"kbswitch/internal/core/switches/models"
	"kbswitch/internal/core/switches/repositories"
)

func New(repo repositories.SwitchesRepo) service {
	return service{repo: repo}
}

type service struct {
	repo repositories.SwitchesRepo
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
