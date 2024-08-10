package models

type SwitchEntity struct {
	Brand            string
	ActuationType    string
	Lifespan         int //in millions
	Name             string
	Image            string
	OperatingForce   int     // in gram-force(gf)
	ActivationTravel float64 // in mm
	TotalTravel      float64 // in mm
	SoundProfile     string  // Quiet, Loud, Normal
	TriggerMethod    string  // mechanical, optical
	Profile          string  // MX, chocov1, chocov2, MX low
}
