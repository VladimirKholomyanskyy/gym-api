package usecase

import (
	"context"
	"fmt"
	"strings"

	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/model"
	"github.com/VladimirKholomyanskyy/gym-api/internal/training/repository"
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
		Description: common.TrimPointer(input.Description),
		ProfileID:   profileID,
	}

	if err := uc.repo.Create(ctx, program); err != nil {
		return nil, fmt.Errorf("failed to create training program: %w", err)
	}
	return program, nil
}

// GetTrainingPrograms retrieves all training programs for the given profile.
func (uc *trainingProgramUseCase) List(ctx context.Context, profileID string, page, pageSize int) ([]model.TrainingProgram, int64, error) {
	if err := common.ValidateUUIDs(profileID); err != nil {
		return nil, 0, err
	}

	programs, totalCount, err := uc.repo.FindByProfileID(ctx, profileID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve training programs: %w", err)
	}
	return programs, totalCount, nil
}

// GetTrainingProgram retrieves a specific training program ensuring it belongs to the profile.
func (uc *trainingProgramUseCase) GetByID(ctx context.Context, profileID, programID string) (*model.TrainingProgram, error) {
	if err := common.ValidateUUIDs(profileID, programID); err != nil {
		return nil, err
	}
	program, err := uc.repo.FindByID(ctx, programID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve training program: %w", err)
	}
	if program == nil {
		return nil, common.NewNotFoundError("training program not found")
	}
	if program.ProfileID != profileID {
		return nil, common.NewForbiddenError("You do not have permission to modify this training program")
	}
	return program, nil
}

// UpdateTrainingProgram updates a training program's details.
func (uc *trainingProgramUseCase) Update(ctx context.Context, profileID, programID string, input openapi.PatchTrainingProgramRequest) (*model.TrainingProgram, error) {
	if err := common.ValidateUUIDs(profileID, programID); err != nil {
		return nil, err
	}
	program, err := uc.GetByID(ctx, profileID, programID)
	if err != nil {
		return nil, err
	}
	updated := false
	if common.HasText(input.Name) {
		program.Name = strings.TrimSpace(*input.Name)
		updated = true
	}
	if common.HasText(input.Description) {
		program.Description = strings.TrimSpace(*input.Description)
		updated = true
	}
	if !updated {
		return program, nil // No changes to apply
	}
	if err = uc.repo.Update(ctx, program); err != nil {
		return nil, fmt.Errorf("failed to update training program: %w", err)
	}
	return program, nil
}

// DeleteTrainingProgram deletes a training program if it belongs to the profile.
func (uc *trainingProgramUseCase) Delete(ctx context.Context, profileID, programID string) error {
	if err := common.ValidateUUIDs(profileID, programID); err != nil {
		return err
	}
	_, err := uc.GetByID(ctx, profileID, programID)
	if err != nil {
		return err
	}
	if err := uc.repo.SoftDelete(ctx, programID, profileID); err != nil {
		return fmt.Errorf("failed to delete training program: %w", err)
	}
	return nil
}
