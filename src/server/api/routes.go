package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	g := engine.Group("v1")

	g.GET("/ping", Ping)
	g.POST("/login", Login)

	us := g.Group("user")

	u := User{}
	us.POST("/register", u.Register)
	us.POST("/unregister", u.Unregister)

	ui := g.Group("view")

	ui.Use(JwtAuthMiddleware())

	vw := View{}
	ui.POST("/enable", vw.Enable)
	ui.POST("/disable", vw.Disable)
	ui.POST("/config", vw.Config)
}
