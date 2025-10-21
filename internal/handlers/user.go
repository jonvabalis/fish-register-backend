package handlers

import (
	"fish-register-backend/internal/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var regData core.UserRegister
	if err := c.ShouldBindJSON(&regData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//TODO validation, hashing, saving, db object
}
