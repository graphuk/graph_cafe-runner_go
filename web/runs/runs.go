package runs

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/cafe-runner/web/runs/models"

	"github.com/labstack/echo"
)

type Handler struct {
	DB *storm.DB
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Handler) Run(c echo.Context) error {
	runID, err := strconv.Atoi(c.Param(`id`))
	check(err)

	model := models.NewRunModel(t.DB, runID)
	return c.Render(http.StatusOK, `runs/views/run.html`, model)
}
