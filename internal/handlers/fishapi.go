package handlers

import "database/sql"

type FishApi struct {
	db *sql.DB
}

func NewFishApi(database *sql.DB) *FishApi {
	return &FishApi{db: database}
}
