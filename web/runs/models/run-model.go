package models

import (
	"github.com/graph-uk/cafe-runner/data/repositories"

	"github.com/asdine/storm"
)

type runModel struct {
	ID       int
	Status   int
	ExitCode string
	StdOut   string
}

func NewRunModel(DB *storm.DB, runID int) *runModel {
	res := &runModel{}

	run := (&repositories.Runs{DB}).Find(runID)
	res.ID = run.ID
	res.Status = int(run.Status)
	res.ExitCode = run.ExitCode
	res.StdOut = string(run.StdOut)

	return res
}
