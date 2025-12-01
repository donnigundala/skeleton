package controllers

import (
	"net/http"
	"strconv"

	"skeleton-v2/app/http/dto"
	"skeleton-v2/app/models"
	"skeleton-v2/app/services"

	"github.com/donnigundala/dg-core/http/request"
	"github.com/donnigundala/dg-core/http/response"
	"github.com/donnigundala/dg-core/validation"
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
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	// Parse and validate JSON
	if err := request.JSONWithValidation(r, &req, c.validator); err != nil {
		if valErr, ok := err.(*validation.Error); ok {
			response.ValidationError(w, valErr.Errors)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}

	// Create user
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := c.service.Create(r.Context(), user); err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	// Return created user
	response.Created(w, c.toResponse(user), "")
}

// Get handles GET /api/v1/users/:id
func (c *UserController) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(request.Param(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	user, err := c.service.GetByID(r.Context(), uint(id))
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	response.Success(w, c.toResponse(user), "User retrieved successfully")
}

// List handles GET /api/v1/users
func (c *UserController) List(w http.ResponseWriter, r *http.Request) {
	page := request.QueryInt(r, "page", 1)
	perPage := request.QueryInt(r, "per_page", 20)

	users, total, err := c.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *c.toResponse(user)
	}

	// Return paginated response
	meta := response.NewPaginationMeta(page, perPage, int(total))
	response.Paginated(w, userResponses, meta)
}

// Update handles PUT /api/v1/users/:id
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(request.Param(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := request.JSONWithValidation(r, &req, c.validator); err != nil {
		if valErr, ok := err.(*validation.Error); ok {
			response.ValidationError(w, valErr.Errors)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}

	// Get existing user
	user, err := c.service.GetByID(r.Context(), uint(id))
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := c.service.Update(r.Context(), user); err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.Success(w, c.toResponse(user), "User updated successfully")
}

// Delete handles DELETE /api/v1/users/:id
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(request.Param(r, "id"), 10, 32)
	if err != nil {
		response.BadRequest(w, "Invalid user ID")
		return
	}

	if err := c.service.Delete(r.Context(), uint(id)); err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}

	response.Success(w, nil, "User deleted successfully")
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
