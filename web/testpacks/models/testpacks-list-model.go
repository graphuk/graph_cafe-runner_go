package models

import (
	"strings"
	"time"

	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"

	"github.com/asdine/storm"
)

type testpackRec struct {
	ID      int
	TimeAgo string
}

type testpacksListModel struct {
	Tesptacks []testpackRec
}

//return duration like 14h27m5s
func timeAgoHumanString(now, moment time.Time) string {
	res := now.Sub(moment).String() //14h27m5.4421878s
	dotIndex := strings.Index(res, `.`)
	if dotIndex != -1 {
		res = res[:dotIndex] + `s` //14h27m5s
	}
	return res
}

func NewTestpacksListModel(DB *storm.DB) *testpacksListModel {
	res := &testpacksListModel{}
	now := time.Now()
	allTestpacks := (&repositories.Testpacks{DB}).FindAll()
	for _, curTestpack := range *allTestpacks {
		res.Tesptacks = append(res.Tesptacks, testpackRec{curTestpack.ID, timeAgoHumanString(now, curTestpack.UploadTime)}) //curTestpack.UploadTime})
	}
	return res
}
