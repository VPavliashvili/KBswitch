package models

type SwitchEntity struct {
	Manufacturer     string
	ActuationType    string
	Lifespan         int //in millions
	Model            string
	Image            []byte
	OperatingForce   int     // in gram-force(gf)
	ActivationTravel float64 // in mm
	TotalTravel      float64 // in mm
	SoundProfile     string  // Quiet, Loud, Normal
	TriggerMethod    string  // mechanical, optical
	Profile          string  // MX, chocov1, chocov2, MX low
}

type Switch struct {
	Brand            string
	ActuationType    string
	Lifespan         int
	Name             string
	Image            string
	OperatingForce   int
	ActivationTravel float64
	TotalTravel      float64
	SoundProfile     string
	TriggerMethod    string
	Profile          string
}
