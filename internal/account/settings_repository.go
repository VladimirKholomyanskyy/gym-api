package account

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// SettingRepository defines the contract for a setting repository.
type SettingRepository interface {
	Create(ctx context.Context, setting *Setting) error
	GetByID(ctx context.Context, id string) (*Setting, error)
	GetByUserID(ctx context.Context, id string) (*Setting, error)
	Update(ctx context.Context, setting *Setting) error
	Delete(ctx context.Context, id string) error
}

// settingRepository is the concrete implementation of SettingRepository.
type settingRepository struct {
	db *gorm.DB
}

// NewSettingRepository returns a new instance of SettingRepository.
func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

// Create inserts a new setting into the database.
func (r *settingRepository) Create(ctx context.Context, setting *Setting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}

// GetByID retrieves a setting by its ID.
func (r *settingRepository) GetByID(ctx context.Context, id string) (*Setting, error) {
	var setting Setting
	err := r.db.WithContext(ctx).First(&setting, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Return nil, nil if no record is found (common Go practice)
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) GetByUserID(ctx context.Context, id string) (*Setting, error) {
	var settings Setting
	err := r.db.WithContext(ctx).Where("profile_id = ?", id).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find settings by profile ID: %w", err)
	}
	return &settings, nil
}

// Update modifies an existing setting.
func (r *settingRepository) Update(ctx context.Context, setting *Setting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

// Delete removes a setting by ID.
func (r *settingRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Setting{}).Error; err != nil {
		return fmt.Errorf("failed to delete user settings: %w", err)
	}
	return nil
}
