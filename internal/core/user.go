package core

import "github.com/gofrs/uuid"

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type UserAuth struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type UserRegister struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
