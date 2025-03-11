package usecase

import (
	"context"
	"strings"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
	"github.com/VladimirKholomyanskyy/gym-api/internal/utils"
)

type TrainingProgramUseCase interface {
	Create(ctx context.Context, profileId string, input openapi.CreateTrainingProgramRequest) (*model.TrainingProgram, error)
	List(ctx context.Context, profileId string, page, pageSize int) ([]model.TrainingProgram, int64, error)
	GetByID(ctx context.Context, profileId, programId string) (*model.TrainingProgram, error)
	Update(ctx context.Context, profileId, programId string, input openapi.PatchTrainingProgramRequest) (*model.TrainingProgram, error)
	Delete(ctx context.Context, profileId, programId string) error
}

type trainingProgramUseCase struct {
	repo repository.TrainingProgramRepository
}

// NewExerciseUseCase creates a new instance of ExerciseUseCase
func NewTrainingProgramUseCase(repo repository.TrainingProgramRepository) TrainingProgramUseCase {
	return &trainingProgramUseCase{
		repo: repo,
	}
}

// CreateTrainingProgram creates a new training program for a given profile.
func (uc *trainingProgramUseCase) Create(ctx context.Context, profileID string, input openapi.CreateTrainingProgramRequest) (*model.TrainingProgram, error) {
	program := &model.TrainingProgram{
		Name:        strings.TrimSpace(input.Name),
		Description: utils.TrimPointer(input.Description),
		ProfileID:   profileID,
	}

	if err := uc.repo.Create(ctx, program); err != nil {
		return nil, err
	}
	return program, nil
}

// GetTrainingPrograms retrieves all training programs for the given profile.
func (uc *trainingProgramUseCase) List(ctx context.Context, profileID string, page, pageSize int) ([]model.TrainingProgram, int64, error) {
	return uc.repo.FindByProfileID(ctx, profileID, page, pageSize)
}

// GetTrainingProgram retrieves a specific training program ensuring it belongs to the profile.
func (uc *trainingProgramUseCase) GetByID(ctx context.Context, profileID, programID string) (*model.TrainingProgram, error) {
	program, err := uc.repo.FindByID(ctx, programID)
	if err != nil {
		return nil, err
	}
	if program.ProfileID != profileID {
		return nil, customerrors.ErrAccessForbidden
	}
	return program, nil
}

// UpdateTrainingProgram updates a training program's details.
func (uc *trainingProgramUseCase) Update(ctx context.Context, profileID, programID string, input openapi.PatchTrainingProgramRequest) (*model.TrainingProgram, error) {
	program, err := uc.GetByID(ctx, profileID, programID)
	if err != nil {
		return nil, err
	}
	updates := make(map[string]any)
	if utils.HasText(input.Name) {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if utils.HasText(input.Description) {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if len(updates) == 0 {
		return program, nil
	}
	return uc.repo.UpdatePartial(ctx, program.ID, updates)
}

// DeleteTrainingProgram deletes a training program if it belongs to the profile.
func (uc *trainingProgramUseCase) Delete(ctx context.Context, profileID, programID string) error {
	_, err := uc.GetByID(ctx, profileID, programID)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, programID)
}
