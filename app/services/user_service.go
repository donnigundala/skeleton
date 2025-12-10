package services

import (
	"context"
	"fmt"
	"time"

	"skeleton-v2/app/models"
	"skeleton-v2/app/repositories"

	cache "github.com/donnigundala/dg-cache"
	"github.com/donnigundala/dg-core/contracts/foundation"
	queue "github.com/donnigundala/dg-queue"
)

// UserService defines the interface for user business logic.
type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context, page, perPage int) ([]*models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
}

// userService implements UserService.
type userService struct {
	repo   repositories.UserRepository
	inject *cache.Injectable
	queue  queue.Queue
}

// NewUserService creates a new user service.
func NewUserService(repo repositories.UserRepository, app foundation.Application, queueManager queue.Queue) UserService {
	return &userService{
		repo:   repo,
		inject: cache.NewInjectable(app),
		queue:  queueManager,
	}
}

// Create creates a new user and dispatches a welcome email job.
func (s *userService) Create(ctx context.Context, user *models.User) error {
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	// Dispatch job to send a welcome email.
	_, err := s.queue.Dispatch("send-welcome-email", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	return err
}

// GetByID retrieves a user by ID with caching.
func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try cache first
	var cachedUser models.User
	err := s.inject.Cache().GetAs(ctx, cacheKey, &cachedUser)
	if err == nil {
		return &cachedUser, nil
	}

	// Cache miss - get from database
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache (5 minutes TTL)
	_ = s.inject.Cache().Put(ctx, cacheKey, user, 5*time.Minute)

	return user, nil
}

// GetAll retrieves all users with pagination.
func (s *userService) GetAll(ctx context.Context, page, perPage int) ([]*models.User, int64, error) {
	return s.repo.GetAll(ctx, page, perPage)
}

// Update updates a user and invalidates cache.
func (s *userService) Update(ctx context.Context, user *models.User) error {
	err := s.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	_ = s.inject.Cache().Forget(ctx, cacheKey)

	return nil
}

// Delete deletes a user and invalidates cache.
func (s *userService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", id)
	_ = s.inject.Cache().Forget(ctx, cacheKey)

	return nil
}
