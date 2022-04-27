package controllers

import (
	"log"
	"net/http"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/gin-gonic/gin"
)

func (app *ProfileBot) UserHandler(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := user.GetByID(userID); err != nil {
		log.Print("Error")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
		return
	}

	c.JSON(200, gin.H{"user": user})
}
