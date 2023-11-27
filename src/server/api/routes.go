package api

import (
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {

	g := engine.Group("v1")

	g.GET("/ping", Ping)
	g.POST("/auth", Auth)
	g.POST("/editor/start", StartEditor)
	g.POST("/editor/stop", StopEditor)
	g.GET("/pods", GetPods)
}
