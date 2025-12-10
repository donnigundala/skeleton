package repositories

import (
	"context"
	"skeleton-v2/app/models"

	"github.com/donnigundala/dg-core/contracts/foundation"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context, page, perPage int) ([]*models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

// userRepository implements UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user.
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID.
func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll retrieves all users with pagination.
func (r *userRepository) GetAll(ctx context.Context, page, perPage int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * perPage
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(perPage).
		Find(&users).Error

	return users, total, err
}

// Update updates a user.
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user by ID.
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// MustResolveUserRepository resolves the user repository from the container.
// It panics if the resolution fails, which is acceptable during app boot.
func MustResolveUserRepository(app foundation.Application) UserRepository {
	repo, err := app.Make("userRepository")
	if err != nil {
		panic("failed to resolve user repository: " + err.Error())
	}
	return repo.(UserRepository)
}
