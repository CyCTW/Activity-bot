package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (app *ProfileBot) CallbackHandler(c *gin.Context) {

	// Early response First
	c.JSON(http.StatusOK, gin.H{"message": "early return"})
	// Verification header
	err := verifyLineSignature(c.Request, app.channel_secret)
	if err != nil {
		log.Print("Verify Failed")
		c.JSON(http.StatusNotFound, gin.H{"message": "Header verification fail"})
		return
	}

	events, err := app.bot.ParseRequest(c.Request)
	if err != nil {
		log.Print("Parse request fail!")
		c.JSON(http.StatusNotFound, gin.H{"message": "Parse request fail"})
		return
	}

	for _, event := range events {
		log.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			err := app.MessageEventHandler(event)
			if err != nil {
				log.Print(err)
			}
		case linebot.EventTypeFollow:
			err := app.WelcomeEventHandler(event)
			if err != nil {
				log.Print(err)
			}
			log.Print("Follow")
		case linebot.EventTypeUnfollow:
			log.Print("Unfollow")
		case linebot.EventTypeJoin:
			err := app.WelcomeEventHandler(event)
			if err != nil {
				log.Print(err)
			}
		case linebot.EventTypeUnsend:
			log.Print("Unsend")
		case linebot.EventTypePostback:
			err := app.PostBackEventHandler(event)
			if err != nil {
				log.Print(err)
			}
		default:
			log.Printf("Unknown event %v", event)
		}

	}

}

func verifyLineSignature(req *http.Request, channel_secret string) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Read error")
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(req.Header.Get("x-line-signature"))
	if err != nil {
		log.Print("decoded error")
		return err
	}
	hash := hmac.New(sha256.New, []byte(channel_secret))
	if _, err := hash.Write(body); err != nil {
		log.Print("Wreite error")
		return err
	}

	res := hmac.Equal(hash.Sum(nil), decoded)
	if !res {
		return errors.New("Verification fail")
	} else {
		log.Print("Verify success!")
	}

	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}
