package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	g := engine.Group("v1")

	g.GET("/ping", Ping)
	g.POST("/login", Login)

	u := User{}
	userGroup := g.Group("user")

	userGroup.POST("/register", u.Register)
	userGroup.POST("/unregister", u.Unregister)

	vw := View{}
	viewsGroup := g.Group("views")

	viewsGroup.Use(JwtAuthMiddleware())

	viewsGroup.GET("", vw.List)
	viewsGroup.GET("/:id", vw.Get)
	viewsGroup.POST("/create", vw.Create)
	viewsGroup.POST("/config", vw.Config)
	viewsGroup.POST("/destroy", vw.Destroy)
}
