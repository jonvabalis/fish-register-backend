package handlers

import (
	"fish-register-backend/internal/core"
	"fish-register-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func (app *FishApi) InsertRod(c *gin.Context) {
	var newRod core.NewRod
	if err := c.ShouldBindJSON(&newRod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUser(c.Request.Context(), app.db, newRod.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	} else if user.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user doesn't exist"})
		return
	}

	rod := core.Rod{
		UUID:          uuid.Must(uuid.NewV4()),
		Nickname:      newRod.Nickname,
		Brand:         newRod.Brand,
		PurchasePlace: newRod.PurchasePlace,
		UserUUID:      newRod.UserUUID,
	}

	if err := db.InsertRod(c.Request.Context(), app.db, rod); err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert rod"})
		return
	}

	c.JSON(201, gin.H{"message": "Rod created"})
}

func (app *FishApi) GetUserRods(c *gin.Context) {
	var req struct {
		UserUUID uuid.UUID `json:"uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rods, err := db.GetUserRods(c.Request.Context(), app.db, req.UserUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch rods"})
		return
	}

	c.JSON(200, gin.H{"rods": rods})
}

func (app *FishApi) PatchRod(c *gin.Context) {
	var rodUpdate core.RodUpdate

	if err := c.ShouldBindJSON(&rodUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rod, err := db.GetRod(c.Request.Context(), app.db, rodUpdate.UUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch rod"})
		return
	}

	rod.ApplyUpdate(rodUpdate)

	if err := db.UpdateRod(c.Request.Context(), app.db, rod); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update rod"})
		return
	}

	c.JSON(200, gin.H{"message": "Rod updated"})
}

func (app *FishApi) DeleteRod(c *gin.Context) {
	var req struct {
		RodUUID uuid.UUID `json:"uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rod, err := db.GetRod(c.Request.Context(), app.db, req.RodUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find rod"})
		return
	} else if rod.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "rod doesn't exist"})
		return
	}

	if err := db.DeleteRod(c.Request.Context(), app.db, rod.UUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rod"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rod deleted successfully"})
}
