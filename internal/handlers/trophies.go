package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (app *FishApi) UploadPicture(c *gin.Context) {
	slot := c.PostForm("slot")

	if slot != "first" && slot != "second" && slot != "third" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot"})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo file required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only jpg format allowed"})
		return
	}

	folder := "./uploads"

	err = os.MkdirAll(folder, 0777)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	dst := fmt.Sprintf("%s/%s.jpg", folder, slot)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded"})
}

func (app *FishApi) DownloadPicture(c *gin.Context) {
	filename := c.Param("filename")

	path := fmt.Sprintf("./uploads/%s", filename)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(path)
}
