package models

import (
	"github.com/graph-uk/cafe-runner/data/repositories"

	"github.com/asdine/storm"
)

type runRec struct {
	ID              int
	DeviceOwnerName string
}

type sessionModel struct {
	ID   int
	Runs []runRec
}

func NewSessionModel(DB *storm.DB, sessionID int) *sessionModel {
	res := &sessionModel{}

	res.ID = sessionID

	allRuns := (&repositories.Runs{DB}).FindBySessionID(sessionID)
	for _, curRun := range *allRuns {
		res.Runs = append(res.Runs, runRec{curRun.ID, curRun.DeviceOwnerName})
	}
	return res
}
