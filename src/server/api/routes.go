package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	v := engine.Group("v1")

	v.GET("/ping", Ping)
	v.POST("/login", Login)

	cd := v.Group("code-editor")

	cd.Use(func(c *gin.Context) {}) // continue

	cd.POST("/create", CreateEditor)
	cd.POST("/delete", DeleteEditor)
}
