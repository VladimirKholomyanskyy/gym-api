package models

type Exercise struct {
	Id               int32    `json:"id"`
	Name             string   `json:"name"`              // Name of the exercise
	Description      string   `json:"description"`       // Optional description of the exercise
	PrimaryMuscle    string   `json:"primary_muscle"`    // Primary muscle targeted
	SecondaryMuscles []string `json:"secondary_muscles"` // List of secondary muscles targeted
	Equipment        string   `json:"equipment"`
}
