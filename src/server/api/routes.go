package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	g := engine.Group("v1")

	g.GET("/ping", Ping)
	g.POST("/login", Login)

	u := User{}
	users := g.Group("users")

	users.GET("", JwtAuthMiddleware(), u.List)
	users.POST("/register", u.Register)
	users.POST("/unregister", u.Unregister)

	vw := View{}
	views := g.Group("views")

	views.Use(JwtAuthMiddleware())

	views.GET("", vw.List)
	views.GET("/:id", vw.Get)
	views.POST("", vw.Create)
	views.PUT("/:id", vw.Config)
	views.DELETE("/:id", vw.Destroy)
}
