package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type ActivityInput struct {
	Name    string `form:"name"`
	Date    string `form:"date"`
	Place   string `form:"place"`
	IdToken string
}

type LineProfileInput struct {
	Sub string
	// Iss     string
	// Aud     string
	// Exp     int
	// Iat     int
	Name string
	// Picture string
}

type LineNotifyResponse struct {
	AccessToken string `json:"access_token"`
}

func getAccessToken(code string, redirect_uri string) (string, error) {

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
	var notifyToken LineNotifyResponse
	log.Print("body", res.Body)
	err = json.NewDecoder(res.Body).Decode(&notifyToken)
	// body, err := ioutil.ReadAll(res.Body)
	// err = json.Unmarshal(body, &notifyToken)
	if err != nil {
		return "", err
	}
	log.Print("ac Token: ", notifyToken.AccessToken)
	return notifyToken.AccessToken, nil
	// Get user data from Line Platform

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

func (app *ProfileBot) NotifyTestGetHandler(c *gin.Context) {
	id := c.Param("id")
	// TODO: Search user id and access token
	var user models.User
	err := user.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"message": "Fail"})
	}
	if err := NotifyUser(user.AccessToken, "GGG"); err != nil {
		c.JSON(400, gin.H{"message": "notify fail"})
	}

	c.JSON(200, gin.H{"message": "success!"})
}

func (app *ProfileBot) NotifyGetHandler(c *gin.Context) {
	code := c.Query("code")
	// State store user_id
	state := c.Query("state")

	user_id := strings.Split(state, "_")[0]
	username := strings.Split(state, "_")[1]

	redirect_uri := os.Getenv("LINE_NOTIFY_REDIRECT_URI")
	log.Print("UID: ", user_id)
	log.Print("code: ", code)
	access_token, err := getAccessToken(code, redirect_uri)
	if err != nil {
		log.Print(err)
		c.JSON(400, gin.H{"message": "失敗!!!"})
		return
	}
	// TODO: Store access token
	var user models.User
	user.StoreAccessToken(user_id, username, access_token)
	c.JSON(200, gin.H{"message": "開啟通知成功!"})

}

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

func (app *ProfileBot) ActivityGetHandler(c *gin.Context) {
	// TODO: Show Single Activity information
	activityID := c.Param("id")
	var activity models.Activity
	if err := activity.GetByID(activityID); err != nil {
		log.Print("Error")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
	}
	users, err := activity.GetInfo()
	if err != nil {
		log.Print("Error")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
	}
	log.Print(users)
	c.JSON(200, gin.H{"activity": activity, "users": users})
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

func (app *ProfileBot) ActivityPostHandler(c *gin.Context) {
	log.Print("Hello!!!!!!!!!!!!!!!!!!!!!")
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
		c.JSON(404, gin.H{"message": "Get profile fail"})
	}

	log.Println("Input: ")
	log.Println(activityInput)
	// Create User if not exists
	user := models.User{LineUserID: lineProfile.Sub, Name: lineProfile.Name}
	log.Print(user)
	if err := user.Create(); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Create User fail"})
		return
	}

	// Create Activity
	dateString := activityInput.Date
	layout := "2006-01-02T15:04:05.000Z"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Date: ", date)

	activity := models.Activity{Name: activityInput.Name, Date: date, Place: activityInput.Place, User: user}
	log.Print(activity)
	if err := activity.Create(); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Create Activity fail"})
		return
	}
	// TODO: Return Activity ID
	c.JSON(200, gin.H{"activityID": activity.ID})
}

func verifyLineSignature(req *http.Request, channel_secret string) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Read error")
	}

	decoded, err := base64.StdEncoding.DecodeString(req.Header.Get("x-line-signature"))
	if err != nil {
		log.Print("decoded error")
	}
	hash := hmac.New(sha256.New, []byte(channel_secret))
	hash.Write(body)

	res := hmac.Equal(hash.Sum(nil), decoded)
	if res != true {
		return errors.New("Verification fail")
	} else {
		log.Print("Verify success!")
	}

	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}

func (app *ProfileBot) CallbackHandler(c *gin.Context) {

	// Early response First
	c.JSON(200, gin.H{"message": "early return"})
	// Verification header
	err := verifyLineSignature(c.Request, app.channel_secret)
	if err != nil {
		log.Print("Verify Failed")
		c.JSON(404, gin.H{"message": "Header verification fail"})
		return
	}

	events, err := app.bot.ParseRequest(c.Request)
	if err != nil {
		log.Print("Parse request fail!")
		c.JSON(404, gin.H{"message": "Parse request fail"})
		return
	}

	for _, event := range events {
		log.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			app.MessageEventHandler(event)
		case linebot.EventTypeFollow:
			log.Print("Follow")
		case linebot.EventTypeUnfollow:
			log.Print("Unfollow")
		case linebot.EventTypePostback:
			app.PostBackEventHandler(event)
		default:
			log.Printf("Unknown event %v", event)
		}

	}

}
