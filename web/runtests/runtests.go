package runtests

import (
	"net/http"

	"github.com/graph-uk/graph_cafe-runner_go/cmd/cafe-runner-server/config"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/web/runtests/models"

	"github.com/labstack/echo"
)

type Handler struct {
	DB               *storm.DB
	CafeRunnerConfig *config.Configuration
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Handler) Runtests(c echo.Context) error {
	model := models.NewRuntestsModel(t.DB, t.CafeRunnerConfig.Server.ExternalURL)
	return c.Render(http.StatusOK, `runtests/views/runtests.html`, model)
}

// func (t *Handler) Runtest(c echo.Context) error {
// 	//testpackID, err := strconv.Atoi(c.Param(`id`))
// 	check(err)

// 	model := models.NewRuntestsModel(t.DB)
// 	return c.Render(http.StatusOK, `runtests/views/testpack.html`, model)
// }
