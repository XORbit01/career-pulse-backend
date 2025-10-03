package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFile godoc
//
//	@Summary		Upload a file (image or document)
//	@Description	Uploads a file and returns its public URL
//	@Tags			Upload
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Param			file	formData	file	true	"File to upload"
//	@Success		200		{object}	models.SuccessResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/upload [post]
func UploadFile(c *gin.Context) {
	// Get configuration from context (set by middleware)
	uploadsPath := c.GetString("uploadsPath")
	staticURL := c.GetString("staticURL")

	if uploadsPath == "" {
		uploadsPath = "./uploads"
	}
	if staticURL == "" {
		staticURL = "/static"
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File upload failed", "error": err.Error()})
		return
	}

	filename := fmt.Sprintf("%s/%d_%s", uploadsPath, time.Now().Unix(), filepath.Base(file.Filename))

	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save file"})
		return
	}

	publicURL := staticURL + "/" + filepath.Base(filename)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File uploaded",
		"data":    map[string]string{"url": publicURL},
	})
}
