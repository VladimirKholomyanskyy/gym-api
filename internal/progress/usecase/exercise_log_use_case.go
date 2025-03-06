package usecase

import (
	"context"
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/repository"
	usecase "github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type LogExerciseUseCase interface {
	Create(ctx context.Context, profileID string, input openapi.CreateExerciseLogRequest) (*model.ExerciseLog, error)
	GetExerciseLog(ctx context.Context, profileID, logID string) (*model.ExerciseLog, error)
	GetExerciseLogsByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetExerciseLogsBySessionID(ctx context.Context, profileID, sessionID string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetExerciseLogsByExerciseID(ctx context.Context, profileID, exerciseId string, page, pageSize int) ([]model.ExerciseLog, int64, error)
	GetWeightPerDay(ctx context.Context, profileID string, exerciseId string, startDate *time.Time, endDate *time.Time) ([]model.WeightPerDay, error)
}
type logExerciseUseCase struct {
	repo    repository.ExerciseLogRepository
	useCase usecase.ExerciseUseCase
}

func NewLogExerciseUseCase(repo repository.ExerciseLogRepository, useCase usecase.ExerciseUseCase) LogExerciseUseCase {
	return &logExerciseUseCase{repo: repo, useCase: useCase}
}
func (uc *logExerciseUseCase) Create(ctx context.Context, profileID string, input openapi.CreateExerciseLogRequest) (*model.ExerciseLog, error) {
	exerciseLog := &model.ExerciseLog{
		SessionID:  input.WorkoutSessionId,
		ExerciseID: input.ExerciseId,
		SetNumber:  int(input.SetNumber),
		Reps:       int(input.RepsCompleted),
		Weight:     float64(input.WeightUsed),
		ProfileID:  profileID}
	err := uc.repo.Create(ctx, exerciseLog)
	if err != nil {
		return nil, err
	}
	return exerciseLog, nil
}

func (uc *logExerciseUseCase) GetExerciseLog(ctx context.Context, profileID, logID string) (*model.ExerciseLog, error) {
	log, err := uc.repo.GetByID(ctx, logID)
	if err != nil {
		return nil, err
	}
	if log.ProfileID != profileID {
		return nil, common.ErrAccessForbidden
	}
	return log, nil

}
func (uc *logExerciseUseCase) GetExerciseLogsByProfileID(ctx context.Context, profileID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	logs, totalCount, err := uc.repo.GetAllByUserID(ctx, profileID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return logs, totalCount, nil
}

func (uc *logExerciseUseCase) GetExerciseLogsBySessionID(ctx context.Context, profileID, sessionID string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	logs, totalCount, err := uc.repo.GetAllByUserIDAndSessionID(ctx, profileID, sessionID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return logs, totalCount, nil
}

func (uc *logExerciseUseCase) GetExerciseLogsByExerciseID(ctx context.Context, profileID, exerciseId string, page, pageSize int) ([]model.ExerciseLog, int64, error) {
	logs, totalCount, err := uc.repo.GetAllByUserIDAndExerciseID(ctx, profileID, exerciseId, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return logs, totalCount, nil
}

func (uc *logExerciseUseCase) GetWeightPerDay(ctx context.Context, profileID string, exerciseId string, startDate *time.Time, endDate *time.Time) ([]model.WeightPerDay, error) {
	weightPerDayArr, err := uc.repo.GetWeightPerDay(ctx, profileID, exerciseId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return weightPerDayArr, err
}
