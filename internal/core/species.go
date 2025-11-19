package core

import "github.com/gofrs/uuid"

type NewSpecies struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type Species struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (s *Species) ApplyUpdate(u Species) {
	if u.Name != "" {
		s.Name = u.Name
	}

	if u.Description != "" {
		s.Description = u.Description
	}
}

func (s *Species) IsEmpty() bool {
	return s.Name == "" && s.Description == ""
}
