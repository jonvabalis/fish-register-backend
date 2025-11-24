package core

import "github.com/gofrs/uuid"

type NewRod struct {
	Nickname      string    `json:"nickname" binding:"required"`
	Brand         string    `json:"brand" binding:"required"`
	PurchasePlace string    `json:"purchasePlace" binding:"required"`
	UserUUID      uuid.UUID `json:"userUUID" binding:"required"`
}

type Rod struct {
	UUID          uuid.UUID `json:"uuid"`
	Nickname      string    `json:"nickname"`
	Brand         string    `json:"brand"`
	PurchasePlace string    `json:"purchasePlace"`
	UserUUID      uuid.UUID `json:"userUUID"`
}

type RodUpdate struct {
	UUID          uuid.UUID `json:"uuid" binding:"required"`
	Nickname      string    `json:"nickname"`
	Brand         string    `json:"brand"`
	PurchasePlace string    `json:"purchasePlace"`
}

func (r *Rod) ApplyUpdate(u RodUpdate) {
	if u.Nickname != "" {
		r.Nickname = u.Nickname
	}

	if u.Brand != "" {
		r.Brand = u.Brand
	}

	if u.PurchasePlace != "" {
		r.PurchasePlace = u.PurchasePlace
	}
}

func (r *Rod) IsEmpty() bool {
	return r.UUID.IsNil() && r.Nickname == "" && r.Brand == "" && r.PurchasePlace == "" && r.UserUUID.IsNil()
}
