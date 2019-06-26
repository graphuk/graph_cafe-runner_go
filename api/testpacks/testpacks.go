package testpacks

import (
	"encoding/base64"
	"log"
	"net/http"
	"strconv"

	"github.com/asdine/storm"

	"github.com/graph-uk/cafe-runner/api/testpacks/models"
	//	"github.com/graph-uk/cafe-runner/data"
	"github.com/graph-uk/cafe-runner/data/repositories"
	"github.com/graph-uk/cafe-runner/logic/services"

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

// Post creates new testpack
func (t *Handler) Post(c echo.Context) error {
	log.Println(`Post testpack received.`)

	model := &testpacks.TestpackPostModel{}
	check(c.Bind(&model))
	testpackContent, err := base64.StdEncoding.DecodeString(model.Content)
	check(err)

	res := (&repositories.Testpacks{t.DB}).Create(&testpackContent)
	log.Println(`Testpack created with id: ` + strconv.Itoa(res.ID))
	go func() { (&services.Testpack{t.DB}).RunInitSteps(res.ID) }() // init testpack async
	return c.JSON(http.StatusCreated, res)
}
