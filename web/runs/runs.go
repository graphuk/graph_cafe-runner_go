package runs

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/web/runs/models"

	"github.com/labstack/echo"
)

type Handler struct {
	DB       *storm.DB
	Hostname *string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Handler) Run(c echo.Context) error {
	runID, err := strconv.Atoi(c.Param(`id`))
	check(err)

	model := models.NewRunModel(t.DB, runID, t.Hostname)
	return c.Render(http.StatusOK, `runs/views/run.html`, model)
}
