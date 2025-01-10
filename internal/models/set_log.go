package models

type SetLog struct {
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	ExerciseLogID uint    `gorm:"not null;index" json:"exercise_log_id"`
	SetNumber     int     `gorm:"not null;check:set_number > 0" json:"set_number"`
	Reps          int     `gorm:"not null;check:reps >= 0" json:"reps"`
	Weight        float64 `gorm:"not null;check:weight >= 0" json:"weight"`

	ExerciseLog ExerciseLog `gorm:"foreignKey:ExerciseLogID;constraint:OnDelete:CASCADE" json:"exercise_log"`
}
