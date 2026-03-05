package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	r := gin.Default()

	r.POST("/register", registerUser)
	r.POST("/login", loginUser)

	r.Run(":8001")
}

func registerUser(c *gin.Context) {

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered",
	})
}

func loginUser(c *gin.Context) {

	var user User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": "dummy-jwt-token",
	})
}