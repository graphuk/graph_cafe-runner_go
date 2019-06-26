package models

import (
	"time"

	"github.com/graph-uk/cafe-runner/data/repositories"

	"github.com/asdine/storm"
)

type sessionRec struct {
	ID          int
	CreatedTime time.Time
}

type testpackModel struct {
	ID       int
	Sessions []sessionRec
}

func NewTestpackModel(DB *storm.DB, testpackID int) *testpackModel {
	res := &testpackModel{}

	res.ID = testpackID

	allSessions := (&repositories.Sessions{DB}).FindByTestpackID(testpackID)
	for _, curSession := range *allSessions {
		res.Sessions = append(res.Sessions, sessionRec{curSession.ID, curSession.CreatedTime})
	}
	return res
}
