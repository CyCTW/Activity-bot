package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/gin-gonic/gin"
)

type LineProfileInput struct {
	Sub string
	// Iss     string
	// Aud     string
	// Exp     int
	// Iat     int
	Name string
	// Picture string
}
type ActivityInput struct {
	Name    string `form:"name"`
	Date    string `form:"date"`
	Place   string `form:"place"`
	IdToken string
}

func (app *ProfileBot) ActivityGetHandler(c *gin.Context) {
	// TODO: Show Single Activity information
	activityID := c.Param("id")
	var activity models.Activity
	if err := activity.GetByID(activityID); err != nil {
		log.Print("Error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Fail"})
		return
	}
	users, err := activity.GetInfo(activityID)
	if err != nil {
		log.Print("Error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Fail"})
		return
	}
	log.Print(users)
	c.JSON(http.StatusOK, gin.H{"activity": activity, "users": users})
}

func (app *ProfileBot) ActivityPostHandler(c *gin.Context) {
	// Bind request
	var activityInput ActivityInput
	if err := c.ShouldBind(&activityInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	idToken := activityInput.IdToken
	log.Print("token", idToken)
	lineProfile, err := getProfile(idToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Get profile fail"})
		return
	}

	log.Println("Input: ")
	log.Println(activityInput)
	// Create User if not exists
	user := models.User{LineUserID: lineProfile.Sub, Name: lineProfile.Name}
	if err := user.Create(); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Create User fail"})
		return
	}
	log.Print("User: <<")
	log.Print(user)
	// Create Activity
	dateString := activityInput.Date
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Time parse fail"})
		return
	}
	fmt.Println("Date: ", date)

	activity := models.Activity{Name: activityInput.Name, Date: date, Place: activityInput.Place, User: user}
	log.Print(activity)
	if err := activity.Create(); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Create Activity fail"})
		return
	}
	// TODO: Return Activity ID
	c.JSON(http.StatusOK, gin.H{"activityName": activity.Name})
}

func getProfile(idToken string) (*LineProfileInput, error) {
	path := "https://api.line.me/oauth2/v2.1/verify"
	postBody := url.Values{}

	postBody.Set("id_token", idToken)
	postBody.Set("client_id", os.Getenv("LINE_LOGIN_CHANNEL_ID"))

	client := &http.Client{}
	r, err := http.NewRequest("POST", path, strings.NewReader(postBody.Encode())) // URL-encoded payload
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	log.Print("Get success!!><")
	lineProfile := LineProfileInput{}

	err = json.NewDecoder(res.Body).Decode(&lineProfile)
	if err != nil {
		return nil, err
	}
	log.Print(lineProfile)

	return &lineProfile, nil

}
