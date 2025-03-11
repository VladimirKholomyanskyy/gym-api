package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/progress/repository"
	training "github.com/VladimirKholomyanskyy/gym-api/internal/training/usecase"
)

type WorkoutSessionUseCase interface {
	StartWorkout(ctx context.Context, profileID string, input openapi.CreateWorkoutSessionRequest) (*model.WorkoutSession, error)
	CompleteWorkout(ctx context.Context, profileID, sessionID string) (*model.WorkoutSession, error)
	List(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error)
	GetByID(ctx context.Context, profileID, sessionID string) (*model.WorkoutSession, error)
}

type workoutSessionUseCase struct {
	repo           repository.WorkoutSessionRepository
	workoutUseCase training.WorkoutUseCase
}

func NewWorkoutSessionUseCase(repo repository.WorkoutSessionRepository,
	workoutUseCase training.WorkoutUseCase) WorkoutSessionUseCase {
	return &workoutSessionUseCase{repo: repo, workoutUseCase: workoutUseCase}
}

func (uc *workoutSessionUseCase) StartWorkout(ctx context.Context, profileID string, input openapi.CreateWorkoutSessionRequest) (*model.WorkoutSession, error) {
	workout, err := uc.workoutUseCase.GetByWorkoutID(ctx, profileID, input.WorkoutId)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(workout)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workout: %w", err)
	}
	workoutSession := &model.WorkoutSession{ProfileID: profileID, WorkoutID: workout.ID, Snapshot: jsonData}
	err = uc.repo.Create(ctx, workoutSession)
	if err != nil {
		return nil, err
	}
	return workoutSession, nil
}

func (uc *workoutSessionUseCase) CompleteWorkout(ctx context.Context, profileID, sessionID string) (*model.WorkoutSession, error) {
	session, err := uc.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.ProfileID != profileID {
		return nil, customerrors.ErrAccessForbidden
	}
	if session.CompletedAt == nil {
		return nil, customerrors.ErrAccessForbidden
	}
	err = uc.repo.UpdatePartial(ctx, session.ID, map[string]any{"completed_at": time.Now()})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *workoutSessionUseCase) List(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error) {
	return uc.repo.GetAllByProfileID(ctx, profileID, page, pageSize)
}

func (uc *workoutSessionUseCase) GetByID(ctx context.Context, profileID, sessionID string) (*model.WorkoutSession, error) {
	session, err := uc.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.ProfileID != profileID {
		return nil, customerrors.ErrAccessForbidden
	}
	return session, err
}
