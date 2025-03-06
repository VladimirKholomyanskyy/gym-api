package account

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ProfileRepository defines CRUD operations for profiles
type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	GetByID(ctx context.Context, id string) (*Profile, error)
	FindByExternalID(ctx context.Context, id string) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	Delete(ctx context.Context, id string) error
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
			return nil, nil // Explicitly return nil for not found cases
		}
		return nil, fmt.Errorf("failed to get profile by ID: %w", err)
	}
	return &profile, nil
}

func (r *profileRepository) FindByExternalID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.db.WithContext(ctx).Where("external_id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find profile by external ID: %w", err)
	}
	return &profile, nil
}

func (r *profileRepository) Update(ctx context.Context, profile *Profile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

func (r *profileRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Profile{}).Error; err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}
