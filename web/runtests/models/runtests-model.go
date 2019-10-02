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

type sessionRec struct {
	ID      int
	TimeAgo string
}

type runtestsModel struct {
	Sessions  []sessionRec
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

func NewRuntestsModel(DB *storm.DB) *runtestsModel {
	res := &runtestsModel{}
	now := time.Now()

	allSessions := (&repositories.Sessions{DB}).FindAllOrderIDDesc()
	for _, curSession := range *allSessions {
		res.Sessions = append(res.Sessions, sessionRec{curSession.ID, timeAgoHumanString(now, curSession.CreatedTime)})
	}

	allTestpacks := (&repositories.Testpacks{DB}).FindAllOrderIDDesc()
	for _, curTestpack := range *allTestpacks {
		res.Tesptacks = append(res.Tesptacks, testpackRec{curTestpack.ID, timeAgoHumanString(now, curTestpack.UploadTime)})
	}

	return res
}
