package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
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
		return nil, common.ErrAccessForbidden
	}
	if session.CompletedAt == nil {
		return nil, common.NewForbiddenError("workout is completed")
	}
	currnetTime := time.Now()
	session.CompletedAt = &currnetTime
	err = uc.repo.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *workoutSessionUseCase) List(ctx context.Context, profileID string, page, pageSize int) ([]model.WorkoutSession, int64, error) {
	sessions, totalCount, err := uc.repo.GetAllByProfileID(ctx, profileID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return sessions, totalCount, err
}

func (uc *workoutSessionUseCase) GetByID(ctx context.Context, profileID, sessionID string) (*model.WorkoutSession, error) {
	session, err := uc.repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.ProfileID != profileID {
		return nil, common.ErrAccessForbidden
	}
	return session, err
}
