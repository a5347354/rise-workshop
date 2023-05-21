package delivery

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func RegisterAggregatorHandler(engine *gin.Engine) error {
	engine.Use(static.Serve("/", static.LocalFile("./web", false)))
	return nil
}
