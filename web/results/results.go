package results

import (
	"net/http"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/web/results/models"

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

func (t *Handler) Results(c echo.Context) error {
	model := models.NewResultsModel(t.DB)
	return c.Render(http.StatusOK, `results/views/results.html`, model)
}
