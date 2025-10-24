package handlers

import (
	"log"
	"net/http"
	"net/mail"

	"fish-register-backend/internal/core"
	"github.com/gin-gonic/gin"
)

func (app *FishApi) Register(c *gin.Context) {
	var regData core.RegisterData
	if err := c.ShouldBindJSON(&regData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Email

	if _, err := mail.ParseAddress(regData.Email); err != nil {
		log.Printf("invalid register email: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Email is in use

	//Username is in use

	//Encrypt password

	//Save user
}
