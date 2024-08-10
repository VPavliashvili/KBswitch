package switches

import (
	"fmt"
	"kbswitch/internal/core/models"
	"strconv"
)

type SwitchDTO struct {
	Brand            string `json:"brand"`
	ActuationType    string `json:"actuationType"`
	Lifespan         string `json:"lifespan"`
	Name             string `json:"name"`
	Image            string `json:"image"`
	OperatingForce   string `json:"operatingForce"`
	ActivationTravel string `json:"activationTravel"`
	TotalTravel      string `json:"totalTravel"`
	SoundProfile     string `json:"SoundProfile"`
	Triggermethod    string `json:"triggermethod"`
	Profile          string `json:"profile"`
}

func asDTO(entity models.SwitchEntity) SwitchDTO {
	lifespan := fmt.Sprintf("%dM", entity.Lifespan)
	opforce := fmt.Sprintf("%dgf", entity.OperatingForce)
	alltravel := strconv.FormatFloat(entity.TotalTravel, 'f', -1, 64) + "mm"
	travel := strconv.FormatFloat(entity.ActivationTravel, 'f', -1, 64) + "mm"

	return SwitchDTO{
		ActivationTravel: travel,
		OperatingForce:   opforce,
		Lifespan:         lifespan,
		TotalTravel:      alltravel,
		Name:             entity.Name,
		Image:            entity.Image,
		Brand:            entity.Brand,
		Profile:          entity.Profile,
		SoundProfile:     entity.SoundProfile,
		Triggermethod:    entity.TriggerMethod,
		ActuationType:    entity.ActuationType,
	}
}
