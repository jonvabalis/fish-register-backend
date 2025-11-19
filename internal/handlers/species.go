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

	species := core.Species{
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
