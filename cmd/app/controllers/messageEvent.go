package controllers

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
		messageArray := strings.SplitN(completeMessage, "-", 2)
		prefixMessage := messageArray[0]
		log.Print(prefixMessage)

		switch prefixMessage {
		case "@我要辦活動":
			if err := app.HandleShowActivityForm(event); err != nil {
				log.Print(err)
				return err
			}
		case "@顯示活動":
			// idx := strings.Index(completeMessage, "ID: ")
			activityName := messageArray[1]
			log.Print(activityName)

			if err := app.HandleShowActivity(event, activityName); err != nil {
				log.Print(err)
				return err
			}
		case "@注意事項":
			err := app.WelcomeEventHandler(event)
			if err != nil {
				log.Print(err)
			}
		}

	default:
		log.Print("Unimplement message event type")
	}
	return nil
}

func (app *ProfileBot) HandleShowActivityForm(event *linebot.Event) error {
	// TODO: Show Form LIFF link
	LIFF_url := os.Getenv("LIFF_URL")
	imageURL := "https://i.imgur.com/Jt8IP8D.jpeg"

	template := linebot.NewButtonsTemplate(
		imageURL, "活動表單", "請點選以下按鈕填寫事件詳細資訊!",
		linebot.NewURIAction("建立活動", LIFF_url),
	)
	if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Buttons alt text", template)).Do(); err != nil {
		return err
	}
	return nil
}

func (app *ProfileBot) HandleShowActivity(event *linebot.Event, activityName string) error {
	// TODO: Perform "attend" button template
	var activity models.Activity
	if err := activity.GetByName(activityName); err != nil {
		text := "查無此活動!"
		if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do(); err != nil {
			return err
		}
		return err
	}

	imageURL := "https://i.imgur.com/HDeWm1D.jpeg"
	actionPayload := fmt.Sprintf("action=attend&activityID=%v", activity.ID)
	title := fmt.Sprintf("活動: %v", activity.Name)

	activity_id := strconv.FormatUint(uint64(activity.ID), 10)
	activityURL := fmt.Sprintf("%v/activity.html?id=%v", os.Getenv("LIFF_URL"), activity_id)

	template := linebot.NewButtonsTemplate(
		imageURL, title, "請點選 \"我要參加\" 按鈕來參加活動，注意要先授權通知後才可參加活動 ",
		// linebot.NewURIAction("授權通知", lineNotifyURL),
		linebot.NewURIAction("查看活動", activityURL),
		linebot.NewPostbackAction("我要參加", actionPayload, "", ""),
	)

	if _, err := app.bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Buttons alt text", template)).Do(); err != nil {
		return err
	}
	return nil
}
