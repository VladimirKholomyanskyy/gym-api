package account

import (
	"context"
	"fmt"

	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"gorm.io/gorm"
)

// SettingRepository defines the contract for a setting repository.
type SettingRepository interface {
	Create(ctx context.Context, setting *Setting) error
	GetByID(ctx context.Context, id string) (*Setting, error)
	GetByProfileID(ctx context.Context, id string) (*Setting, error)
	Update(ctx context.Context, setting *Setting) error
	Delete(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
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
	if err := r.db.WithContext(ctx).Create(setting).Error; err != nil {
		return fmt.Errorf("failed to create settings: %w", err)
	}
	return nil
}

// GetByID retrieves a setting by its ID.
func (r *settingRepository) GetByID(ctx context.Context, id string) (*Setting, error) {
	var setting Setting
	err := r.db.WithContext(ctx).First(&setting, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotFound
		}
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) GetByProfileID(ctx context.Context, id string) (*Setting, error) {
	var settings Setting
	err := r.db.WithContext(ctx).Where("profile_id = ?", id).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to find settings by profile ID: %w", err)
	}
	return &settings, nil
}

// Update modifies an existing setting.
func (r *settingRepository) Update(ctx context.Context, setting *Setting) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(setting).Error; err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
		return nil
	})
}

// Delete removes a setting by ID.
func (r *settingRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Setting{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete settings: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrEntityNotFound
	}
	return nil
}

func (r *settingRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&Setting{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete settings: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrEntityNotFound
	}
	return nil
}
