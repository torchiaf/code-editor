package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	g := engine.Group("v1")

	g.GET("/ping", Ping)
	g.POST("/login", Login)

	ui := g.Group("ui")

	ui.Use(JwtAuthMiddleware())

	// ui.GET("/user",controllers.CurrentUser)
	ui.POST("/create", Create)
	ui.POST("/delete", Delete)
}
