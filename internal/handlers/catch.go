package handlers

import (
	"fish-register-backend/internal/core"
	"fish-register-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"time"
)

func (app *FishApi) CreateCatch(c *gin.Context) {
	var catchData core.CreateCatchData
	if err := c.ShouldBindJSON(&catchData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUser(c.Request.Context(), app.db, catchData.UsersUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	} else if user.IsEmpty() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user doesn't exist"})
		return
	}

	if catchData.SpeciesUUID != nil && !catchData.SpeciesUUID.IsNil() {
		species, err := db.GetSpecies(c.Request.Context(), app.db, *catchData.SpeciesUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create catch"})
			return
		} else if species.IsEmpty() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "species does not exist"})
			return
		}
	}

	if catchData.LocationsUUID != nil && !catchData.LocationsUUID.IsNil() {
		location, err := db.GetLocation(c.Request.Context(), app.db, *catchData.LocationsUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create catch"})
			return
		} else if location.IsEmpty() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "location does not exist"})
			return
		}
	}

	if catchData.RodsUUID != nil && !catchData.RodsUUID.IsNil() {
		rod, err := db.GetRod(c.Request.Context(), app.db, *catchData.RodsUUID)
		if err != nil {
			log.Printf("error getting rod: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create catch"})
			return
		} else if rod.IsEmpty() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rod does not exist"})
			return
		}

		if !rod.UserUUID.IsNil() && rod.UserUUID != catchData.UsersUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "rod does not belong to user"})
			return
		}
	}

	if catchData.Length != nil && *catchData.Length < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "length cannot be negative"})
		return
	}

	if catchData.Weight != nil && *catchData.Weight < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weight cannot be negative"})
		return
	}

	// Create catch object
	catch := core.Catch{
		UUID:          uuid.Must(uuid.NewV4()),
		Nickname:      catchData.Nickname,
		Length:        catchData.Length,
		Weight:        catchData.Weight,
		Comment:       catchData.Comment,
		CaughtAt:      catchData.CaughtAt,
		CreatedAt:     time.Now().UTC(),
		SpeciesUUID:   catchData.SpeciesUUID,
		LocationsUUID: catchData.LocationsUUID,
		UsersUUID:     catchData.UsersUUID,
		RodsUUID:      catchData.RodsUUID,
	}

	if err := db.CreateCatch(c.Request.Context(), app.db, catch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create catch"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"catch": catch})
}

func (app *FishApi) GetUserCatches(c *gin.Context) {
	userUUIDStr := c.Param("userUUID")
	if userUUIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userUUID is required"})
		return
	}

	userUUID, err := uuid.FromString(userUUIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userUUID format"})
		return
	}

	user, err := db.GetUser(c.Request.Context(), app.db, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch catches"})
		return
	} else if user.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	catches, err := db.GetUserCatches(c.Request.Context(), app.db, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch catches"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"catches": catches})
}

func (app *FishApi) UpdateUserCatch(c *gin.Context) {
	var updateData core.UpdateCatchData
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.UUID.IsNil() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "catch uuid is required"})
		return
	}

	existingCatch, err := db.GetCatch(c.Request.Context(), app.db, updateData.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update catch"})
		return
	} else if existingCatch.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "catch not found"})
		return
	}

	if updateData.SpeciesUUID != nil && !updateData.SpeciesUUID.IsNil() {
		species, err := db.GetSpecies(c.Request.Context(), app.db, *updateData.SpeciesUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update catch"})
			return
		} else if species.IsEmpty() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "species does not exist"})
			return
		}
	}

	if updateData.LocationsUUID != nil && !updateData.LocationsUUID.IsNil() {
		location, err := db.GetLocation(c.Request.Context(), app.db, *updateData.LocationsUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update catch"})
			return
		} else if location.IsEmpty() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "location does not exist"})
			return
		}
	}

	if updateData.RodsUUID != nil && !updateData.RodsUUID.IsNil() {
		rod, err := db.GetRod(c.Request.Context(), app.db, *updateData.RodsUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update catch"})
			return
		} else if rod.IsEmpty() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rod does not exist"})
			return
		}

		if !rod.UserUUID.IsNil() && rod.UserUUID != existingCatch.UsersUUID {
			c.JSON(http.StatusForbidden, gin.H{"error": "rod does not belong to catch owner"})
			return
		}
	}

	if updateData.Length != nil && *updateData.Length < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "length cannot be negative"})
		return
	}
	if updateData.Weight != nil && *updateData.Weight < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weight cannot be negative"})
		return
	}

	if err := db.UpdateCatch(c.Request.Context(), app.db, updateData.UUID, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update catch"})
		return
	}

	updatedCatch, err := db.GetCatch(c.Request.Context(), app.db, updateData.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve updated catch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"catch": updatedCatch})
}

func (app *FishApi) DeleteCatch(c *gin.Context) {
	var req struct {
		CatchUUID uuid.UUID `json:"uuid" binding:"required"`
		UserUUID  uuid.UUID `json:"user_uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	catch, err := db.GetCatch(c.Request.Context(), app.db, req.CatchUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete catch"})
		return
	} else if catch.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "catch not found"})
		return
	}

	if catch.UsersUUID != req.UserUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to delete this catch"})
		return
	}

	if err := db.DeleteCatch(c.Request.Context(), app.db, req.CatchUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete catch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "catch deleted successfully"})
}
