package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	model "github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func HasText(s *string) bool {
	return s != nil && strings.TrimSpace(*s) != ""
}

func TrimPointer(s *string) string {
	if s == nil {
		return ""
	}
	trimmed := strings.TrimSpace(*s)
	return trimmed
}

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func ParseTime(dateStr string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func CalculateTotalPages(totalRecords int64, pageSize int32) int32 {
	if pageSize <= 0 {
		return 0
	}
	return int32((totalRecords + int64(pageSize) - 1) / int64(pageSize))
}

func ValidateUUIDs(ids ...string) error {
	for _, id := range ids {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("%w: %s", customerrors.ErrInvalidUUID, id)
		}
	}
	return nil
}

func ErrorResponse(status int, code openapi.ErrorCodes, message string) (openapi.ImplResponse, error) {
	return openapi.Response(status, openapi.ErrorResponse{ErrorCode: code, Message: message}), nil
}

func ConvertExercise(gormExercise *model.Exercise) *openapi.Exercise {
	return &openapi.Exercise{
		Id:              gormExercise.ID,
		Name:            gormExercise.Name,
		PrimaryMuscle:   gormExercise.PrimaryMuscle,
		SecondaryMuscle: gormExercise.SecondaryMuscle,
		Equipment:       gormExercise.Equipment,
		Description:     gormExercise.Description,
	}
}

func ConvertExercises(gormExercises []model.Exercise) []openapi.Exercise {
	apiExercises := make([]openapi.Exercise, len(gormExercises))
	for i, e := range gormExercises {
		apiExercises[i] = *ConvertExercise(&e)
	}
	return apiExercises
}

func ConvertTrainingProgram(gromProgram *model.TrainingProgram) *openapi.TrainingProgram {
	return &openapi.TrainingProgram{
		Id:          gromProgram.ID,
		Name:        gromProgram.Name,
		Description: gromProgram.Description,
	}

}

func ConvertTrainingPrograms(gromPrograms []model.TrainingProgram) []openapi.TrainingProgram {
	apiTrainingPrograms := make([]openapi.TrainingProgram, len(gromPrograms))
	for i, e := range gromPrograms {
		apiTrainingPrograms[i] = *ConvertTrainingProgram(&e)
	}
	return apiTrainingPrograms
}

func ConvertWorkout(gormWorkout *model.Workout) *openapi.Workout {
	return &openapi.Workout{
		Id:       gormWorkout.ID,
		Name:     gormWorkout.Name,
		Position: int32(gormWorkout.Position),
	}
}

func ConvertWorkouts(gormWorkouts []model.Workout) []openapi.Workout {
	apiWorkouts := make([]openapi.Workout, len(gormWorkouts))
	for i, e := range gormWorkouts {
		apiWorkouts[i] = *ConvertWorkout(&e)
	}
	return apiWorkouts
}

func ConvertWorkoutExercise(gormWorkoutExercise *model.WorkoutExercise) *openapi.WorkoutExercise {
	return &openapi.WorkoutExercise{
		Id:         gormWorkoutExercise.ID,
		WorkoutId:  gormWorkoutExercise.WorkoutID,
		ExerciseId: gormWorkoutExercise.ExerciseID,
		Sets:       int32(gormWorkoutExercise.Sets),
		Reps:       int32(gormWorkoutExercise.Reps),
		Position:   int32(gormWorkoutExercise.Position),
	}
}

func ConvertWorkoutExercises(gormWorkoutExercise []model.WorkoutExercise) []openapi.WorkoutExercise {
	apiWorkoutExercises := make([]openapi.WorkoutExercise, len(gormWorkoutExercise))
	for i, e := range gormWorkoutExercise {
		apiWorkoutExercises[i] = *ConvertWorkoutExercise(&e)
	}
	return apiWorkoutExercises
}

func ConvertScheduledWorkout(gormScheduledWorkout *model.ScheduledWorkout) *openapi.ScheduledWorkout {
	return &openapi.ScheduledWorkout{
		Id:        gormScheduledWorkout.ID,
		WorkoutId: gormScheduledWorkout.WorkoutID,
		Date:      gormScheduledWorkout.Date.Format("2006-01-02"),
		Notes:     gormScheduledWorkout.Notes,
	}
}

func ConvertScheduledWorkouts(gormScheduledWorkouts []model.ScheduledWorkout) []openapi.ScheduledWorkout {
	apiScheduledWorkouts := make([]openapi.ScheduledWorkout, len(gormScheduledWorkouts))
	for i, e := range gormScheduledWorkouts {
		apiScheduledWorkouts[i] = *ConvertScheduledWorkout(&e)
	}
	return apiScheduledWorkouts
}

func ConverWorkoutSnapshot(workoutJson *datatypes.JSON) (*openapi.WorkoutSessionWorkoutSnapshot, error) {
	var workout model.Workout
	err := json.Unmarshal(*workoutJson, &workout)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workout JSON: %w", err)
	}
	var snapshotExercises []openapi.WorkoutExercise
	for _, we := range workout.Exercises {
		snapshotExercises = append(snapshotExercises, openapi.WorkoutExercise{
			Id:         we.ID,
			WorkoutId:  we.WorkoutID,
			ExerciseId: we.ExerciseID,
			Position:   int32(we.Position),
			Sets:       int32(we.Sets),
			Reps:       int32(we.Reps),
		})
	}
	return &openapi.WorkoutSessionWorkoutSnapshot{
		Id:                workout.ID,
		Name:              workout.Name,
		TrainingProgramId: workout.TrainingProgramID,
		WorkoutExercises:  snapshotExercises,
	}, nil
}
