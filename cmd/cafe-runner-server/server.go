package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gobuffalo/packr"
	"github.com/graph-uk/graph_cafe-runner_go/api/runs"
	"github.com/graph-uk/graph_cafe-runner_go/api/sessions"
	"github.com/graph-uk/graph_cafe-runner_go/cmd/cafe-runner-server/config"

	//"github.com/labstack/echo/middleware"

	"github.com/graph-uk/graph_cafe-runner_go/api/testpacks"
	"github.com/graph-uk/graph_cafe-runner_go/data/models"
	webhome "github.com/graph-uk/graph_cafe-runner_go/web/home"
	webresults "github.com/graph-uk/graph_cafe-runner_go/web/results"
	webruns "github.com/graph-uk/graph_cafe-runner_go/web/runs"
	webruntests "github.com/graph-uk/graph_cafe-runner_go/web/runtests"
	websessions "github.com/graph-uk/graph_cafe-runner_go/web/sessions"
	webtestbylink "github.com/graph-uk/graph_cafe-runner_go/web/testbylink"
	webtestpacks "github.com/graph-uk/graph_cafe-runner_go/web/testpacks"

	"github.com/labstack/echo"
)

type CafeRunnerServer struct {
	config *config.Configuration
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *CafeRunnerServer) initFolders(datapath string) {
	check(os.Chdir(datapath))
	check(os.MkdirAll(`_data`, 0644))
	check(os.MkdirAll(`_data`+string(os.PathSeparator)+`testpacks`, 0644))
	check(os.MkdirAll(`_data`+string(os.PathSeparator)+`runs`, 0644))
}

func (t *CafeRunnerServer) openDB(path string) *storm.DB {
	db, err := storm.Open(path)
	check(err)
	return db
}

func (t *CafeRunnerServer) initDB(db *storm.DB) {
	tx, err := db.Begin(true)
	check(err)
	defer tx.Rollback()

	check(tx.Init(&models.Testpack{}))
	check(tx.ReIndex(&models.Testpack{}))
	check(tx.Init(&models.Session{}))
	check(tx.ReIndex(&models.Session{}))
	check(tx.Init(&models.Run{}))
	check(tx.ReIndex(&models.Run{}))

	check(tx.Commit())
}

// Start web server
func (t *CafeRunnerServer) Start() {
	t.initFolders(t.config.Data.Path)

	db := t.openDB(`_data/base.db`)
	defer db.Close()
	t.initDB(db)
	assetsBox := packr.NewBox("../../web")

	templates, _ := parseTemplates(&assetsBox)
	renderer := &Template{Templates: templates}

	e := echo.New()
	e.Renderer = renderer
	//e.Use(middleware.Logger())

	//assets
	assetHandler := http.FileServer(assetsBox)
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))

	//web
	e.GET("/", (&webhome.Handler{db}).Home)
	e.GET("/testpacks", (&webtestpacks.Handler{db}).TestpacksList)
	e.GET("/testpacks/:id", (&webtestpacks.Handler{db}).Testpack)
	e.GET("/sessions/:id", (&websessions.Handler{db}).Session)
	e.GET("/runs/:id", (&webruns.Handler{db, &t.config.Server.Hostname}).Run)
	e.GET("/runtests", (&webruntests.Handler{db, t.config}).Runtests)
	e.GET("/results", (&webresults.Handler{db}).Results)
	e.GET("/testbylink", (&webtestbylink.Handler{}).TestByLink)

	//api
	e.POST("/api/v1/testpacks", (&testpacks.Handler{db}).Post)
	e.POST("/api/v1/sessions", (&sessions.Handler{db}).Post)
	e.POST("/api/v1/runs", (&runs.Handler{db, t.config}).Post)
	e.GET("/api/v1/runs/:id", (&runs.Handler{db, t.config}).Get)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(t.config.Server.Port)))
}
