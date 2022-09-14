package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kh9543/bot/apps/bot"
	"github.com/kh9543/bot/apps/http/heartbeat"
)

func main() {
	if err := bot.NewDiscordBot("!", "token").Start(); err != nil {
		panic(err)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiv1 := r.Group("/api/v1")

	heartbeat.NewHandler(apiv1)
	r.Run(":80")
}
