package home

import (
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/cafe-runner/data/repositories"
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

func (t *Handler) Home(c echo.Context) error {
	lastSession := (&repositories.Sessions{t.DB}).FindLast()
	c.Response().Header().Set(`Cache-Control`, `no-cache`)
	c.Response().Header().Set(`Pragma`, `no-cache`)
	return c.Redirect(http.StatusMovedPermanently, `/sessions/`+strconv.Itoa(lastSession.ID))
}
