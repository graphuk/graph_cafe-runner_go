package results

import (
	"net/http"
	//	"strconv"

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

// func (t *Handler) Runtest(c echo.Context) error {
// 	//testpackID, err := strconv.Atoi(c.Param(`id`))
// 	check(err)

// 	model := models.NewResultsModel(t.DB)
// 	return c.Render(http.StatusOK, `results/views/testpack.html`, model)
// }
