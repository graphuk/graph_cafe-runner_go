package models

import (
	"strings"
	"time"

	//	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"

	"github.com/asdine/storm"
)

type runtestRec struct {
	ID      int
	TimeAgo string
}

type runtestsListModel struct {
	Tesptacks []runtestRec
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

func NewRuntestsListModel(DB *storm.DB) *runtestsListModel {
	res := &runtestsListModel{}
	//now := time.Now()
	// allRuntests := (&repositories.Runtests{DB}).FindAll()
	// for _, curRuntest := range *allRuntests {
	// 	res.Tesptacks = append(res.Tesptacks, runtestRec{curRuntest.ID, timeAgoHumanString(now, curRuntest.UploadTime)})
	// }
	return res
}
