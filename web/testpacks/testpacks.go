package testpacks

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/cafe-runner/web/testpacks/models"

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

func (t *Handler) TestpacksList(c echo.Context) error {
	model := models.NewTestpacksListModel(t.DB)
	return c.Render(http.StatusOK, `testpacks/views/testpacks-list.html`, model)
}

func (t *Handler) Testpack(c echo.Context) error {
	testpackID, err := strconv.Atoi(c.Param(`id`))
	check(err)

	model := models.NewTestpackModel(t.DB, testpackID)
	return c.Render(http.StatusOK, `testpacks/views/testpack.html`, model)
}
