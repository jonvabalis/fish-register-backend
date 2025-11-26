package core

import (
	"github.com/gofrs/uuid"
	"time"
)

type Catch struct {
	UUID          uuid.UUID  `json:"uuid"`
	Nickname      string     `json:"nickname"`
	Length        *float64   `json:"length"`
	Weight        *float64   `json:"weight"`
	Comment       string     `json:"comment"`
	CaughtAt      time.Time  `json:"caught_at"`
	CreatedAt     time.Time  `json:"created_at"`
	SpeciesUUID   *uuid.UUID `json:"species_uuid"`
	LocationsUUID *uuid.UUID `json:"locations_uuid"`
	UsersUUID     uuid.UUID  `json:"users_uuid"`
	RodsUUID      *uuid.UUID `json:"rods_uuid"`
}

type CreateCatchData struct {
	Nickname      string     `json:"nickname"`
	Length        *float64   `json:"length"`
	Weight        *float64   `json:"weight"`
	Comment       string     `json:"comment"`
	CaughtAt      time.Time  `json:"caught_at" binding:"required"`
	SpeciesUUID   *uuid.UUID `json:"species_uuid"`
	LocationsUUID *uuid.UUID `json:"locations_uuid"`
	UsersUUID     uuid.UUID  `json:"users_uuid" binding:"required"`
	RodsUUID      *uuid.UUID `json:"rods_uuid"`
}

type UpdateCatchData struct {
	UUID          uuid.UUID  `json:"uuid" binding:"required"`
	Nickname      *string    `json:"nickname"`
	Length        *float64   `json:"length"`
	Weight        *float64   `json:"weight"`
	Comment       *string    `json:"comment"`
	CaughtAt      *time.Time `json:"caught_at"`
	SpeciesUUID   *uuid.UUID `json:"species_uuid"`
	LocationsUUID *uuid.UUID `json:"locations_uuid"`
	RodsUUID      *uuid.UUID `json:"rods_uuid"`
}

func (c *Catch) IsEmpty() bool {
	return c.UUID.IsNil() && c.UsersUUID.IsNil()
}
