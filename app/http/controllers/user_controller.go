package controllers

import (
	"net/http"
	"strconv"

	"skeleton-v2/app/http/dto"
	"skeleton-v2/app/models"
	"skeleton-v2/app/services"

	"github.com/donnigundala/dg-core/validation"
	"github.com/gin-gonic/gin"
)

// UserController handles user HTTP requests.
type UserController struct {
	service   services.UserService
	validator *validation.Validator
}

// NewUserController creates a new user controller.
func NewUserController(service services.UserService, validator *validation.Validator) *UserController {
	return &UserController{
		service:   service,
		validator: validator,
	}
}

// Create handles POST /api/v1/users
func (c *UserController) Create(ctx *gin.Context) {
	var req dto.CreateUserRequest

	// Bind JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate struct
	if err := c.validator.ValidateStruct(ctx.Request.Context(), &req); err != nil {
		if valErr, ok := err.(*validation.Error); ok {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": valErr.Errors})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := c.service.Create(ctx.Request.Context(), user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return created user
	ctx.JSON(http.StatusCreated, c.toResponse(user))
}

// Get handles GET /api/v1/users/:id
func (c *UserController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.service.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, c.toResponse(user))
}

// List handles GET /api/v1/users
func (c *UserController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	users, total, err := c.service.GetAll(ctx.Request.Context(), page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTOs
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *c.toResponse(user)
	}

	// Return paginated response
	ctx.JSON(http.StatusOK, gin.H{
		"data": userResponses,
		"meta": gin.H{
			"current_page": page,
			"per_page":     perPage,
			"total":        total,
		},
	})
}

// Update handles PUT /api/v1/users/:id
func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validator.ValidateStruct(ctx.Request.Context(), &req); err != nil {
		if valErr, ok := err.(*validation.Error); ok {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": valErr.Errors})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing user
	user, err := c.service.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := c.service.Update(ctx.Request.Context(), user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, c.toResponse(user))
}

// Delete handles DELETE /api/v1/users/:id
func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// toResponse converts a model to a response DTO.
func (c *UserController) toResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
