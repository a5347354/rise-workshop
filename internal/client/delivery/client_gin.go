package delivery

import (
	"github.com/a5347354/rise-workshop/internal/client"
	"github.com/a5347354/rise-workshop/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type clientHandler struct {
	usecase client.Usecase
}

func RegisterClientHandler(engine *gin.Engine,
	usecase client.Usecase) error {
	h := &clientHandler{
		usecase,
	}
	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/event", h.create)
		}
	}
	return nil
}

type request struct {
}

func (handler clientHandler) create(c *gin.Context) {
	var req request
	if err := c.ShouldBind(&req); err != nil {
		pkg.AbortWithErrorJSON(c, http.StatusBadRequest, err)
		return
	}
	err := handler.usecase.SendMessage(c)
	if err != nil {
		logrus.Warn(err)
	}

	c.Status(http.StatusNoContent)
}
