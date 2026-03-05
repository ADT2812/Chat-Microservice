package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var onlineUsers = make(map[string]bool)

func main() {

	r := gin.Default()

	r.POST("/online/:user", setOnline)
	r.POST("/offline/:user", setOffline)
	r.GET("/status/:user", getStatus)

	r.Run(":8004")
}

func setOnline(c *gin.Context) {

	user := c.Param("user")
	onlineUsers[user] = true

	c.JSON(http.StatusOK, gin.H{"status": "online"})
}

func setOffline(c *gin.Context) {

	user := c.Param("user")
	delete(onlineUsers, user)

	c.JSON(http.StatusOK, gin.H{"status": "offline"})
}

func getStatus(c *gin.Context) {

	user := c.Param("user")

	if onlineUsers[user] {
		c.JSON(http.StatusOK, gin.H{"online": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"online": false})
	}
}
