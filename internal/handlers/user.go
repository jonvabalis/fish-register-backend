package handlers

import (
	"fish-register-backend/internal/db"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
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

	if _, err := mail.ParseAddress(regData.Email); err != nil {
		log.Printf("invalid register email: %s", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user, err := db.GetUserByEmail(c.Request.Context(), app.db, regData.Email); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	} else if !user.IsEmpty() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check email"})
		return
	}

	if user, err := db.GetUserByUsername(c.Request.Context(), app.db, regData.Username); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
		return
	} else if !user.IsEmpty() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check username"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	user := core.UserAuth{
		UUID:     uuid.Must(uuid.NewV4()),
		Username: regData.Username,
		Email:    regData.Email,
		Password: string(hashedPassword),
	}

	if err := db.CreateUser(c.Request.Context(), app.db, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
