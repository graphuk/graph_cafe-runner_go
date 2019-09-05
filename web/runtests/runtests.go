package runtests

import (
	"net/http"
	//	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/web/runtests/models"

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

func (t *Handler) RuntestsList(c echo.Context) error {
	model := models.NewRuntestsListModel(t.DB)
	return c.Render(http.StatusOK, `runtests/views/runtests.html`, model)
}

// func (t *Handler) Runtest(c echo.Context) error {
// 	//testpackID, err := strconv.Atoi(c.Param(`id`))
// 	check(err)

// 	model := models.NewRuntestsListModel(t.DB)
// 	return c.Render(http.StatusOK, `runtests/views/testpack.html`, model)
// }
