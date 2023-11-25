package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	g := engine.Group("v1")

	g.GET("/ping", Ping)
	g.POST("/editor", CreateEditor)
}
