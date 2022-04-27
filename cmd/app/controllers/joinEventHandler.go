package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (app *ProfileBot) WelcomeEventHandler(event *linebot.Event) error {
	welcomeMessage_1 := "歡迎使用此機器人! 使用前請先\n1. 將此機器人設為好友\n2. 點選下方授權通知"
	welcomeMessage_2 := "指令:\n輸入 \"@我要辦活動\" 來顯示辦活動的表單\n輸入 \"@顯示活動-[活動名稱]\" 來查看活動詳細資訊\n感謝你的使用><"
	imageURL := "https://i.imgur.com/Jt8IP8D.jpeg"
	notifyURL := fmt.Sprintf("%v/auth.html", os.Getenv("LIFF_URL"))
	var messages []linebot.SendingMessage
	messages = append(messages, linebot.NewTextMessage(welcomeMessage_1))
	messages = append(messages, linebot.NewTemplateMessage("Buttons alt text",
		linebot.NewButtonsTemplate(imageURL, "授權通知!", "請點選以下按鈕授權通知",
			linebot.NewURIAction("授權通知", notifyURL))))
	messages = append(messages, linebot.NewTextMessage(welcomeMessage_2))

	if _, err := app.bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
