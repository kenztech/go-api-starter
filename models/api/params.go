package api

import "time"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserRequest struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=admin merchant operator"`
	Status   string `json:"status" validate:"required,oneof=active inactive banned"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

type UsersResponse struct {
	Success    bool       `json:"success"`
	Users      []UserData `json:"users"`
	Pagination Pagination `json:"pagination"`
}

type UserResponse struct {
	Success bool     `json:"success"`
	User    UserData `json:"user"`
}

type UserData struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
