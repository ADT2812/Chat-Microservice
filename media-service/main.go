package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	os.MkdirAll("./uploads", os.ModePerm)

	r := gin.Default()

	r.POST("/upload", uploadFile)

	r.Run(":8003")
}

func uploadFile(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
		return
	}

	path := "./uploads/" + file.Filename

	c.SaveUploadedFile(file, path)

	c.JSON(http.StatusOK, gin.H{
		"file": file.Filename,
	})
}
