package models

type SwitchRequestBody struct {
	Brand            string  `json:"brand"`
	ActuationType    string  `json:"actuationType"`
	Lifespan         int     `json:"lifespan"`
	Name             string  `json:"name"`
	Image            string  `json:"image"`
	OperatingForce   int     `json:"operatingForce"`
	ActivationTravel float64 `json:"activationTravel"`
	TotalTravel      float64 `json:"totalTravel"`
	SoundProfile     string  `json:"soundProfile"`
	TriggerMethod    string  `json:"triggerMethod"`
	Profile          string  `json:"profile"`
}
