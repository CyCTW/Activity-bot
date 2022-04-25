package controllers

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (app *ProfileBot) MessageEventHandler(event *linebot.Event) error {
	// Do something
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		completeMessage := message.Text
		log.Print(completeMessage)
		messageArray := strings.Split(completeMessage, ", ")
		prefixMessage := messageArray[0]
		log.Print(prefixMessage)

		switch prefixMessage {
		case "顯示活動表單":
			if err := app.HandleShowActivity(event); err != nil {
				log.Print(err)
				return err
			}
		case "我要舉辦活動":
			// idx := strings.Index(completeMessage, "ID: ")
			activityID := messageArray[1][4:]
			log.Print(activityID)

			if err := app.HandleCreateActivity(event, activityID); err != nil {
				log.Print(err)
				return err
			}
		}
		// if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
		// 	log.Print(err)
		// 	return err
		// }
	default:
		log.Print("Unimplement message event type")
	}
	return nil
}

func (app *ProfileBot) HandleShowActivity(event *linebot.Event) error {
	// TODO: Show Form LIFF link
	LIFF_url := os.Getenv("LIFF_URL")
	imageURL := "https://image.com"

	template := linebot.NewButtonsTemplate(
		imageURL, "Build your activity!", "請點選以下按鈕填寫事件詳細資訊!",
		linebot.NewURIAction("建立事件", LIFF_url),
	)
	if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Buttons alt text", template)).Do(); err != nil {
		return err
	}
	return nil
}

func (app *ProfileBot) HandleCreateActivity(event *linebot.Event, activityID string) error {
	// TODO: Perform "attend" button template
	imageURL := "https://image.com"
	actionPayload := fmt.Sprintf("action=attend&activityID=%v", activityID)

	var activity models.Activity
	if err := activity.GetByID(activityID); err != nil {
		return err
	}
	redirect_uri := os.Getenv("LINE_NOTIFY_REDIRECT_URI")
	lineNotifyClientID := os.Getenv("LINE_NOTIFY_CLIENT_ID")

	title := fmt.Sprintf("活動: %v", activity.Name)
	userID := event.Source.UserID
	profile, err := app.bot.GetProfile(userID).Do()
	if err != nil {
		return err
	}
	state := fmt.Sprintf("%v_%v", userID, profile.DisplayName)

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", lineNotifyClientID)
	params.Add("redirect_uri", redirect_uri)
	params.Add("scope", "notify")
	params.Add("state", state)
	lineNotifyURL := fmt.Sprintf("https://notify-bot.line.me/oauth/authorize?%v", params.Encode())

	template := linebot.NewButtonsTemplate(
		imageURL, title, "若要參加，請先點選\"授權通知\"，再按下我要\"參加按鈕\"",
		linebot.NewURIAction("授權通知", lineNotifyURL),
		linebot.NewPostbackAction("我要參加", actionPayload, "", ""),
	)

	if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Buttons alt text", template)).Do(); err != nil {
		return err
	}
	return nil
}
