package sessions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/asdine/storm"

	"github.com/graph-uk/graph_cafe-runner_go/api/sessions/models"
	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"

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

// Post creates new session
func (t *Handler) Post(c echo.Context) error {
	log.Println(`Post session received.`)

	model := &sessions.SessionPostModel{}

	if err := c.Bind(&model); err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	res := (&repositories.Sessions{t.DB}).Create()

	log.Println(`Session is created with id: ` + strconv.Itoa(res.ID))
	return c.JSON(http.StatusCreated, res)
}
