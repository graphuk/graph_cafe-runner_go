package sessions

import (
	"log"
	"net/http"
	"strconv"

	"github.com/asdine/storm"

	"github.com/graph-uk/graph_cafe-runner_go/api/sessions/models"
	"github.com/graph-uk/graph_cafe-runner_go/logic/services"

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

	res := (&services.Session{t.DB}).Create(model.TestpackID)
	log.Println(`Session for testpack ` + strconv.Itoa(model.TestpackID) + ` created with id: ` + strconv.Itoa(res.ID))
	return c.JSON(http.StatusCreated, res)
}

//// Post creates new session
//func Get(c echo.Context) error {
//	log.Println(`Get_recvd`)
//	return c.Redirect(301, `http://localhost:21000/browser/connect`)
//}
