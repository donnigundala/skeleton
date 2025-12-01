package dto

// CreateUserRequest represents the request to create a user.
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email,max=100"`
}

// UpdateUserRequest represents the request to update a user.
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=3,max=100"`
	Email string `json:"email" validate:"omitempty,email,max=100"`
}

// UserResponse represents a user in API responses.
type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
