package models

import (
	"time"

	"github.com/graph-uk/cafe-runner/data/repositories"

	"github.com/asdine/storm"
)

type testpackRec struct {
	ID         int
	UploadTime time.Time
}

type testpacksListModel struct {
	Tesptacks []testpackRec
}

func NewTestpacksListModel(DB *storm.DB) *testpacksListModel {
	res := &testpacksListModel{}
	allTestpacks := (&repositories.Testpacks{DB}).FindAll()
	for _, curTestpack := range *allTestpacks {
		res.Tesptacks = append(res.Tesptacks, testpackRec{curTestpack.ID, curTestpack.UploadTime})
	}
	return res
}
