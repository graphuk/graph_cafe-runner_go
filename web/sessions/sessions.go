package sessions

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/cafe-runner/web/sessions/models"

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

func (t *Handler) Session(c echo.Context) error {
	sessionID, err := strconv.Atoi(c.Param(`id`))
	check(err)

	model := models.NewSessionModel(t.DB, sessionID)
	return c.Render(http.StatusOK, `sessions/views/session.html`, model)
}
