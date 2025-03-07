package account

import (
	"context"
	"fmt"

	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"gorm.io/gorm"
)

// ProfileRepository defines CRUD operations for profiles
type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	BatchCreate(ctx context.Context, profiles []*Profile) error
	GetByID(ctx context.Context, id string) (*Profile, error)
	GetByIDWithAssociations(ctx context.Context, id string, preload ...string) (*Profile, error)
	FindByExternalID(ctx context.Context, id string) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	UpdatePartial(ctx context.Context, id string, updates map[string]any) error
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	PermanentDelete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Profile, int64, error)
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

func (r *profileRepository) BatchCreate(ctx context.Context, profiles []*Profile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, profile := range profiles {
			if err := tx.Create(profile).Error; err != nil {
				return fmt.Errorf("failed to create profile: %w", err)
			}
		}
		return nil
	})
}

func (r *profileRepository) GetByID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) GetByIDWithAssociations(ctx context.Context, id string, preload ...string) (*Profile, error) {
	var profile Profile
	query := r.db.WithContext(ctx)

	for _, relation := range preload {
		query = query.Preload(relation)
	}

	err := query.Where("id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) FindByExternalID(ctx context.Context, id string) (*Profile, error) {
	var profile Profile
	err := r.db.WithContext(ctx).Where("external_id = ?", id).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrEntityNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) Update(ctx context.Context, profile *Profile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(profile).Error; err != nil {
			return fmt.Errorf("failed to update profile: %w", err)
		}
		return nil
	})
}

func (r *profileRepository) UpdatePartial(ctx context.Context, id string, updates map[string]any) error {
	result := r.db.WithContext(ctx).Model(&Profile{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrEntityNotFound
	}
	return nil
}

func (r *profileRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Profile{})
	if result.Error != nil {
		return fmt.Errorf("failed to create profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrEntityNotFound
	}
	return nil
}

func (r *profileRepository) Restore(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&Profile{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to restore profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrEntityNotFound
	}
	return nil
}

func (r *profileRepository) PermanentDelete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&Profile{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return common.ErrEntityNotFound
	}
	return nil
}

func (r *profileRepository) List(ctx context.Context, limit, offset int) ([]*Profile, int64, error) {
	var profiles []*Profile
	var count int64

	err := r.db.WithContext(ctx).Model(&Profile{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&profiles).Error
	if err != nil {
		return nil, 0, err
	}

	return profiles, count, nil
}
