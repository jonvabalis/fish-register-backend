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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	} else if !user.IsEmpty() {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}

	if user, err := db.GetUserByUsername(c.Request.Context(), app.db, regData.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	} else if !user.IsEmpty() {
		c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
		return
	}

	if len(regData.Password) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is too short"})
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

func (app *FishApi) GetUsers(c *gin.Context) {
	users, err := db.GetUsers(c.Request.Context(), app.db)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(200, gin.H{"users": users})
}

func (app *FishApi) ChangeLogin(c *gin.Context) {
	var authData core.UserAuth
	if err := c.ShouldBindJSON(&authData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authData.UUID.IsNil() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patch object provided"})
		return
	}

	user, err := db.GetUser(c.Request.Context(), app.db, authData.UUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user doesnt exist"})
		return
	}

	if len(authData.Email) > 0 {
		if _, err := mail.ParseAddress(authData.Email); err != nil {
			log.Printf("invalid email: %s", err)

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if foundUser, err := db.GetUserByEmail(c.Request.Context(), app.db, authData.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		} else if !foundUser.IsEmpty() {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}

		user.Email = authData.Email
	}

	if len(authData.Username) > 0 {
		if foundUser, err := db.GetUserByUsername(c.Request.Context(), app.db, authData.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		} else if !foundUser.IsEmpty() {
			c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
			return
		}

		user.Username = authData.Username
	}

	if len(authData.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(authData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		user.Password = string(hashedPassword)
	}

	if err := db.UpdateUser(c.Request.Context(), app.db, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (app *FishApi) DeleteUser(c *gin.Context) {
	var req struct {
		UserUUID uuid.UUID `json:"uuid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUser(c.Request.Context(), app.db, req.UserUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	} else if user.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user doesn't exist"})
		return
	}

	if err := db.DeleteUser(c.Request.Context(), app.db, user.UUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
