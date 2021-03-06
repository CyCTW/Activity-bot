package controllers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (app *ProfileBot) PostBackEventHandler(event *linebot.Event) error {
	// Do something
	// data format: action=[attend]&activityName=[Name]
	log.Print("Postback Event: ", event.Postback)

	data := event.Postback.Data
	splitArray := strings.Split(data, "&")
	if len(splitArray) == 0 {
		return errors.New("Wrong Postback Data format")
	}
	log.Print(splitArray)

	action := strings.Split(splitArray[0], "=")[1]
	activityID := strings.Split(splitArray[1], "=")[1]

	// Get user profile
	userID := event.Source.UserID
	profile, err := app.bot.GetProfile(userID).Do()

	if err != nil {
		log.Print("User profile not found")
		return err
	}
	switch action {
	case "attend":

		var user models.User
		err := user.GetByLineID(userID)
		if err != nil {
			text := "使用者未開啟通知，請先開啟通知後再嘗試!"
			if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
				log.Print(err)
				return err
			}
			return err
		}
		log.Print(profile)
		var activity models.Activity
		err = activity.GetByID(activityID)
		if err != nil {
			return err
		}
		// Add user to attendee
		if err := activity.AddParticipants(&user); err != nil {
			return err
		}
		message := fmt.Sprintf("%v 參加活動 \"%v\" 成功!", user.Name, activity.Name)
		if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
			return err
		}
		// Start scheduler
		log.Print("Start scheduler")
		StartScheduler(&user, &activity)
	default:
		log.Print("Unimplement postback event type")
	}

	return nil
}
