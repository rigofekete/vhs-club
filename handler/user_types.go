package handler

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=4,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUserBatchRequest struct {
	// dive validator needed to iterate each object of the slice and apply the CreateUserRequest binded validations
	Users []CreateUserRequest `json:"users" binding:"required,dive"`
}

type UserResponse struct {
	PublicID uuid.UUID `json:"public_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserLoginResponse struct {
	PublicID uuid.UUID `json:"public_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Token    string    `json:"token"`
}

type UserBatchResponse struct {
	Users         []UserResponse `json:"users"`
	AlreadyExists int32          `json:"already_exists"`
}
