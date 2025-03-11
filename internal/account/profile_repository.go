package account

import (
	"context"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"gorm.io/gorm"
)

// ProfileRepository defines CRUD operations for profiles
type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	GetByID(ctx context.Context, id string) (*Profile, error)
	GetByExternalID(ctx context.Context, id string) (*Profile, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) error
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
}

// profileRepository implements ProfileRepository using GORM
type profileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new instance of ProfileRepository
func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) Create(ctx context.Context, profile *Profile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return nil
}

func (r *profileRepository) GetByID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch profile by id: %w", err)
	}
	return &profile, nil
}

func (r *profileRepository) GetByExternalID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.db.WithContext(ctx).Where("external_id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch profile by external id: %w", err)
	}
	return &profile, nil
}

func (r *profileRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&Profile{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

func (r *profileRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Profile{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

func (r *profileRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&Profile{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}
