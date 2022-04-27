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
	/* Debug Mode */
	// r := gin.Default()

	/* Production Mode */
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// r.LoadHTMLGlob("./views/*")
	r.Static("assets", "../nextjs/out")
	r.POST("/callback", app.CallbackHandler)
	r.GET("/activity/:id", app.ActivityGetHandler)
	r.GET("/user/:id", app.UserHandler)
	r.GET("/notify", app.NotifyGetHandler)
	r.GET("/notify_test/:id", app.NotifyTestGetHandler)
	r.POST("/activity", app.ActivityPostHandler)
	r.Run()

}
