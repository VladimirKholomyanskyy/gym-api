package handlers

import (
	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
)

func convertExerciseLog(gormExerciseLog *model.ExerciseLog) *openapi.ExerciseLog {
	return &openapi.ExerciseLog{
		Id:               gormExerciseLog.ID,
		ExerciseId:       gormExerciseLog.ExerciseID,
		WorkoutSessionId: gormExerciseLog.SessionID,
		SetNumber:        int32(gormExerciseLog.SetNumber),
		RepsCompleted:    int32(gormExerciseLog.Reps),
		WeightUsed:       int32(gormExerciseLog.Weight),
		LoggedAt:         gormExerciseLog.LoggedAt,
	}

}

func convertExerciseLogs(gormExerciseLogs []model.ExerciseLog) []openapi.ExerciseLog {
	apiExerciseLogs := make([]openapi.ExerciseLog, len(gormExerciseLogs))
	for i, e := range gormExerciseLogs {
		apiExerciseLogs[i] = *convertExerciseLog(&e)
	}
	return apiExerciseLogs
}
