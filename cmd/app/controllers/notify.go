package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/gin-gonic/gin"
)

type LineNotifyResponse struct {
	AccessToken string `json:"access_token"`
}

func (app *ProfileBot) NotifyTestGetHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Search user id and access token
	var user models.User
	err := user.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Fail"})
	}
	if err := NotifyUser(user.AccessToken, "test"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "notify fail"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "success!"})
}

func (app *ProfileBot) NotifyGetHandler(c *gin.Context) {
	code := c.Query("code")
	// State store user_id
	state := c.Query("state")

	user_id := strings.Split(state, "_")[0]
	username := strings.Split(state, "_")[1]

	log.Print("UID: ", user_id)
	log.Print("code: ", code)
	access_token, err := getAccessToken(code)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "失敗，請稍後再試"})
		return
	}
	// Store access token
	var user models.User
	if err := user.StoreAccessToken(user_id, username, access_token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "存取notify token錯誤"})
	}
	c.HTML(http.StatusOK, "notify.html", gin.H{"title": "website"})

}

func NotifyUser(access_token string, message string) error {
	log.Print("Enter Notify!")
	path := "https://notify-api.line.me/api/notify"

	// access_token := user.AccessToken
	log.Print(access_token)
	auth_header := fmt.Sprintf("Bearer %v", access_token)
	postBody := url.Values{}
	postBody.Set("message", message)

	client := &http.Client{}
	r, err := http.NewRequest("POST", path, strings.NewReader(postBody.Encode())) // URL-encoded payload
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", auth_header)

	_, err = client.Do(r)
	if err != nil {
		return err
	}
	return nil
}

func getAccessToken(code string) (string, error) {
	redirect_uri := os.Getenv("LINE_NOTIFY_REDIRECT_URI")
	client_id := os.Getenv("LINE_NOTIFY_CLIENT_ID")
	client_secret := os.Getenv("LINE_NOTIFY_CLIENT_SECRET")

	path := "https://notify-bot.line.me/oauth/token"
	postBody := url.Values{}

	postBody.Set("grant_type", "authorization_code")
	postBody.Set("code", code)
	postBody.Set("redirect_uri", redirect_uri)
	postBody.Set("client_secret", client_secret)
	postBody.Set("client_id", client_id)

	client := &http.Client{}
	r, err := http.NewRequest("POST", path, strings.NewReader(postBody.Encode())) // URL-encoded payload
	if err != nil {
		return "", err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	log.Print("Get success!!><")
	log.Print(res.Status)
	if res.StatusCode != http.StatusOK {
		return "", errors.New("Get access token failed")
	}

	var notifyToken LineNotifyResponse
	log.Print("body", res.Body)
	err = json.NewDecoder(res.Body).Decode(&notifyToken)
	if err != nil {
		return "", err
	}
	log.Print("ac Token: ", notifyToken.AccessToken)
	return notifyToken.AccessToken, nil
	// Get user data from Line Platform

}
