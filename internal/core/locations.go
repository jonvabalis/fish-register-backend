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

func (l *Location) ApplyUpdate(u Location) {
	if u.Name != "" {
		l.Name = u.Name
	}
	if u.Address != "" {
		l.Address = u.Address
	}
	if u.Type != "" {
		l.Type = u.Type
	}
}

func (l *Location) IsEmpty() bool {
	return l.Name == "" && l.Address == "" && l.Type == ""
}
