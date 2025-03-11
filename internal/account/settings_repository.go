package account

import (
	"context"
	"fmt"

	customerrors "github.com/VladimirKholomyanskyy/gym-api/internal/customErrors"
	"gorm.io/gorm"
)

// SettingRepository defines the contract for a setting repository.
type SettingRepository interface {
	Create(ctx context.Context, setting *Setting) error
	GetByID(ctx context.Context, id string) (*Setting, error)
	GetByProfileID(ctx context.Context, id string) (*Setting, error)
	UpdatePartial(ctx context.Context, id string, updates map[string]any) error
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
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to fetch by settings by id: %w", err)
	}
	return &setting, nil
}

func (r *settingRepository) GetByProfileID(ctx context.Context, id string) (*Setting, error) {
	var settings Setting
	err := r.db.WithContext(ctx).Where("profile_id = ?", id).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerrors.ErrEntityNotFound
		}
		return nil, fmt.Errorf("failed to find settings by profile ID: %w", err)
	}
	return &settings, nil
}

func (r *settingRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&Setting{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update settings: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

// Delete removes a setting by ID.
func (r *settingRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Setting{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete settings: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}

func (r *settingRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&Setting{})
	if result.Error != nil {
		return fmt.Errorf("failed to permanent delete settings: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return customerrors.ErrEntityNotFound
	}
	return nil
}
