package testbylink

import (
	"log"
	"net/http"

	"github.com/graph-uk/graph_cafe-runner_go/web/testbylink/models"
	"github.com/labstack/echo"
)

type Handler struct {
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Handler) TestByLink(c echo.Context) error {
	model := models.NewRuntestsModel(c.QueryParam("session"), c.QueryParam("testpack"), c.QueryParam("name"))
	log.Println(c.Render(http.StatusOK, `testbylink/views/testbylink.html`, model))
	return c.Render(http.StatusOK, `testbylink/views/testbylink.html`, model)
}
