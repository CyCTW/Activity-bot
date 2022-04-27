package controllers

import (
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (app *ProfileBot) WelcomeEventHandler(event *linebot.Event) error {
	welcomeMessage_1 := "歡迎使用此機器人! \n使用前請先將此機器人設為好友，並且開啟通知"
	welcomeMessage_2 := "指令:\n輸入 \"@我要辦活動\" 來顯示辦活動的表單\n輸入 \"@顯示活動-[活動名稱]\" 來查看活動詳細資訊\n感謝你的使用><"
	var messages []linebot.SendingMessage
	messages = append(messages, linebot.NewTextMessage(welcomeMessage_1))
	messages = append(messages, linebot.NewTextMessage(welcomeMessage_2))

	if _, err := app.bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
		log.Print(err)
		return err
	}
	return nil
}
