package home

import (
	"net/http"

	"github.com/asdine/storm"
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
	c.Response().Header().Set(`Cache-Control`, `no-cache`)
	c.Response().Header().Set(`Pragma`, `no-cache`)
	return c.Redirect(http.StatusMovedPermanently, `/runtests`)
}
