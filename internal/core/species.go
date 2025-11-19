package core

import "github.com/gofrs/uuid"

type NewSpecies struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"type" binding:"required"`
}

type Species struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (l *Species) ApplyUpdate(u Species) {
	if u.Name == "" {
		l.Name = u.Name
	}

	if u.Description != "" {
		l.Description = u.Description
	}
}
