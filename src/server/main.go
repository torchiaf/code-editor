package main

import (
	"github.com/torchiaf/code-editor/server/api"

	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()
	engine.Use(api.CORSMiddleware())

	api.Routes(engine)

	engine.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
