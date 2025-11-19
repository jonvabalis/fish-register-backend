package handlers

import (
	"fish-register-backend/internal/core"
	"fish-register-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func (app *FishApi) InsertLocation(c *gin.Context) {
	var newLocation core.NewLocation
	if err := c.ShouldBindJSON(&newLocation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location := core.Location{
		UUID:    uuid.Must(uuid.NewV4()),
		Name:    newLocation.Name,
		Address: newLocation.Address,
		Type:    newLocation.Type,
	}

	if err := db.InsertLocation(c.Request.Context(), app.db, location); err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert location"})
		return
	}

	c.JSON(201, gin.H{"message": "Location created"})
}

func (app *FishApi) GetLocations(c *gin.Context) {
	locations, err := db.GetLocations(c.Request.Context(), app.db)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch locations"})
		return
	}

	c.JSON(200, gin.H{"locations": locations})
}

func (app *FishApi) PatchLocation(c *gin.Context) {
	var locationUpdate core.Location
	if err := c.ShouldBindJSON(&locationUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := db.GetLocation(c.Request.Context(), app.db, locationUpdate.UUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch location"})
		return
	}

	location.ApplyUpdate(locationUpdate)

	if err := db.UpdateLocation(c.Request.Context(), app.db, location); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update location"})
		return
	}

	c.JSON(200, gin.H{"message": "Location updated"})
}
