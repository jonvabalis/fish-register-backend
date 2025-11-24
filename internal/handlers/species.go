package handlers

import (
	"fish-register-backend/internal/core"
	"fish-register-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func (app *FishApi) InsertSpecies(c *gin.Context) {
	var newSpecies core.NewSpecies
	if err := c.ShouldBindJSON(&newSpecies); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	species, err := db.GetSpeciesByName(c.Request.Context(), app.db, newSpecies.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch species"})
		return
	} else if !species.IsEmpty() {
		c.JSON(200, gin.H{"message": "Species already exists"})
		return
	}

	species = core.Species{
		UUID:        uuid.Must(uuid.NewV4()),
		Name:        newSpecies.Name,
		Description: newSpecies.Description,
	}

	if err := db.InsertSpecies(c.Request.Context(), app.db, species); err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert species"})
		return
	}

	c.JSON(201, gin.H{"message": "Species created"})
}

func (app *FishApi) GetAllSpecies(c *gin.Context) {
	species, err := db.GetMultipleSpecies(c.Request.Context(), app.db)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch species"})
		return
	}

	c.JSON(200, gin.H{"species": species})
}

func (app *FishApi) PatchSpecies(c *gin.Context) {
	var speciesUpdate core.Species
	if err := c.ShouldBindJSON(&speciesUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	species, err := db.GetSpecies(c.Request.Context(), app.db, speciesUpdate.UUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch species"})
		return
	}

	species.ApplyUpdate(speciesUpdate)

	if err := db.UpdateSpecies(c.Request.Context(), app.db, species); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update species"})
		return
	}

	c.JSON(200, gin.H{"message": "Species updated"})
}

func (app *FishApi) DeleteSpecies(c *gin.Context) {
	var req struct {
		SpeciesUUID uuid.UUID `json:"uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	species, err := db.GetSpecies(c.Request.Context(), app.db, req.SpeciesUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find species"})
		return
	} else if species.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "species doesn't exist"})
		return
	}

	if err := db.DeleteSpecies(c.Request.Context(), app.db, species.UUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete species"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Species deleted successfully"})
}

func (app *FishApi) GetAllSpeciesByLocation(c *gin.Context) {
	var req struct {
		UUID uuid.UUID `json:"locationUUID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := db.GetLocation(c.Request.Context(), app.db, req.UUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch location"})
		return
	} else if location.IsEmpty() {
		c.JSON(404, gin.H{"error": "Failed to find provided location"})
		return
	}

	species, err := db.GetMultipleSpeciesByLocation(c.Request.Context(), app.db, location.UUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch species"})
		return
	}

	c.JSON(200, gin.H{"species": species})
}

func (app *FishApi) InsertSpeciesToLocation(c *gin.Context) {
	var req struct {
		LocationUUID uuid.UUID `json:"locationUUID" binding:"required"`
		SpeciesUUID  uuid.UUID `json:"speciesUUID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := db.GetLocation(c.Request.Context(), app.db, req.LocationUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch location"})
		return
	} else if location.IsEmpty() {
		c.JSON(404, gin.H{"error": "Failed to find the provided location"})
		return
	}

	species, err := db.GetSpecies(c.Request.Context(), app.db, req.SpeciesUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch species"})
		return
	} else if species.IsEmpty() {
		c.JSON(404, gin.H{"error": "Failed to find the provided species"})
		return
	}

	if err := db.InsertSpeciesToLocation(c.Request.Context(), app.db, species.UUID, location.UUID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert species to location"})
		return
	}

	c.JSON(201, gin.H{"message": "Species added to location"})
}
func (app *FishApi) DeleteSpeciesFromLocation(c *gin.Context) {
	var req struct {
		LocationUUID uuid.UUID `json:"locationUUID" binding:"required"`
		SpeciesUUID  uuid.UUID `json:"speciesUUID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location, err := db.GetLocation(c.Request.Context(), app.db, req.LocationUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch location"})
		return
	} else if location.IsEmpty() {
		c.JSON(404, gin.H{"error": "Failed to find the provided location"})
		return
	}

	species, err := db.GetSpecies(c.Request.Context(), app.db, req.SpeciesUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch species"})
		return
	} else if species.IsEmpty() {
		c.JSON(404, gin.H{"error": "Failed to find the provided species"})
		return
	}

	if err := db.DeleteSpeciesFromLocation(c.Request.Context(), app.db, species.UUID, location.UUID); err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete species from location"})
		return
	}

	c.JSON(200, gin.H{"message": "Species deleted from location successfully"})
}
