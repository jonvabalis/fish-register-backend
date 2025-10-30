package core

import "github.com/gofrs/uuid"

type NewLocation struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type Location struct {
	UUID    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Type    string    `json:"type"`
}
