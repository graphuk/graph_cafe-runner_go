package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gobuffalo/packr"
	"github.com/graph-uk/cafe-runner/api/runs"
	"github.com/graph-uk/cafe-runner/api/sessions"

	//"github.com/labstack/echo/middleware"

	"github.com/graph-uk/cafe-runner/api/testpacks"
	"github.com/graph-uk/cafe-runner/data/models"
	webhome "github.com/graph-uk/cafe-runner/web/home"
	webruns "github.com/graph-uk/cafe-runner/web/runs"
	websessions "github.com/graph-uk/cafe-runner/web/sessions"
	webtestpacks "github.com/graph-uk/cafe-runner/web/testpacks"

	"github.com/labstack/echo"
)

type CafeRunnerServer struct {
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *CafeRunnerServer) initFolders() {
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
func (t *CafeRunnerServer) Start(port int) {
	t.initFolders()

	db := t.openDB(`_data/base.db`)
	defer db.Close()
	t.initDB(db)
	assetsBox := packr.NewBox("../../web")

	templates, _ := parseTemplates(&assetsBox)
	renderer := &Template{Templates: templates}

	//assets
	e := echo.New()
	e.Renderer = renderer
	//e.Use(middleware.Logger())
	assetHandler := http.FileServer(assetsBox)
	e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))

	//web
	e.GET("/", (&webhome.Handler{db}).Home)
	e.GET("/testpacks", (&webtestpacks.Handler{db}).TestpacksList)
	e.GET("/testpacks/:id", (&webtestpacks.Handler{db}).Testpack)
	e.GET("/sessions/:id", (&websessions.Handler{db}).Session)
	e.GET("/runs/:id", (&webruns.Handler{db}).Run)

	//api
	e.POST("/api/v1/testpacks", (&testpacks.Handler{db}).Post)
	e.POST("/api/v1/sessions", (&sessions.Handler{db}).Post)
	e.POST("/api/v1/runs", (&runs.Handler{db}).Post)
	e.GET("/api/v1/runs/:id", (&runs.Handler{db}).Get)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}
