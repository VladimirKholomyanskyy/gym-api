package models

type WorkoutExercise struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	WorkoutID  uint    `gorm:"not null;index" json:"workout_id"`  // Foreign key to workouts
	ExerciseID uint    `gorm:"not null;index" json:"exercise_id"` // Foreign key to exercises
	Sets       int     `gorm:"not null;check:sets > 0" json:"sets"`
	Reps       int     `gorm:"not null;check:reps > 0" json:"reps"`
	Weight     float64 `gorm:"not null;check:weight >= 0" json:"weight"`
	Workout    Workout `gorm:"foreignKey:WorkoutID"` // Association with Workout
}
