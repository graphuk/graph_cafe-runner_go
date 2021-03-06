package runs

import (
	"log"
	"net/http"
	"strconv"

	//	"time"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/api/runs/models"
	"github.com/graph-uk/graph_cafe-runner_go/cmd/cafe-runner-server/config"
	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"
	"github.com/graph-uk/graph_cafe-runner_go/logic/services"

	"github.com/labstack/echo"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	DB               *storm.DB
	CafeRunnerConfig *config.Configuration
}

// Post creates new run
func (t *Handler) Post(c echo.Context) error {
	log.Println(`Post run received.`)

	model := &runs.RunPostModel{}
	if err := c.Bind(&model); err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	sessionID, err := strconv.Atoi(model.SessionID)
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	testpackID, err := strconv.Atoi(model.TestpackID)
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	res := (&services.Run{t.DB, t.CafeRunnerConfig}).Create(sessionID, testpackID, model.DeviceOwnerName)
	log.Println(`Run for session ` + model.SessionID + ` created with id: ` + strconv.Itoa(res.ID))
	go func() {
		(&services.Run{t.DB, t.CafeRunnerConfig}).RunInitSteps(res.ID, t.CafeRunnerConfig.Server.ExternalURL)
	}() // init testpack async
	return c.Redirect(301, `/api/v1/runs/`+strconv.Itoa(res.ID))
}

func (t *Handler) Get(c echo.Context) error {
	log.Println(`Get run received.`)
	runID, err := strconv.Atoi(c.Param(`id`))
	check(err)

	//	time.Sleep(15 * time.Second)

	run := (&repositories.Runs{t.DB}).Find(runID)
	// if run == nil {
	// 	return c.JSON(http.StatusNotFound, `nil`)
	// }
	//return c.Redirect(301, `http://localhost:`+strconv.Itoa(run.Port)+`/browser/connect`)
	return c.JSON(http.StatusOK, run)
}
