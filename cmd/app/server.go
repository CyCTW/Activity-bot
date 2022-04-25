package main

import (
	"log"

	"github.com/cyctw/line-profile-bot/cmd/app/controllers"
	"github.com/cyctw/line-profile-bot/cmd/app/models"
	"github.com/gin-gonic/gin"
)

func main() {

	app, err := controllers.Init()
	if err != nil {
		log.Fatal(err)
	}
	models.ConnectDatabase()

	r := gin.Default()
	r.Static("assets", "./static/line-liff-v2-starter/src/nextjs/out")
	r.POST("/callback", app.CallbackHandler)
	// r.POST("/notify", app.NotifyPostHandler)
	r.GET("/activity/:id", app.ActivityGetHandler)
	r.GET("/user/:id", app.UserHandler)
	r.GET("/notify", app.NotifyGetHandler)
	r.GET("/notify_test/:id", app.NotifyTestGetHandler)
	r.POST("/activity", app.ActivityPostHandler)
	r.Run()

}
