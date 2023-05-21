package delivery

import (
	"github.com/a5347354/rise-workshop/internal/aggregator"
	"github.com/a5347354/rise-workshop/pkg"

	"fmt"
	"net/http"
	
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type aggregatorHandler struct {
	usecase aggregator.Usecase
}

func RegisterAggregatorHandler(engine *gin.Engine, usecase aggregator.Usecase) error {
	engine.Use(static.Serve("/", static.LocalFile("./web", false)))
	h := &aggregatorHandler{usecase: usecase}
	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/events", h.List)
		}
	}
	return nil
}

func (handler aggregatorHandler) List(c *gin.Context) {
	events, err := handler.usecase.ListEventByKeyword(c, c.Query("keyword"))
	if err != nil {
		logrus.Error(err)
		pkg.AbortWithErrorJSON(c, http.StatusInternalServerError, fmt.Errorf("something wrong"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
