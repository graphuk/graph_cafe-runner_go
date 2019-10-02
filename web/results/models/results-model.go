package models

import (
	"log"
	"strings"
	"time"

	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"

	"github.com/asdine/storm"
)

const (
	ResultStatusNotRunned      = 0
	ResultStatusLastPassed     = 1
	ResultStatusLastLastFailed = 2
)

type resultRec struct {
	Status     byte
	TriesCount int
	LastTryAgo string
}

type sessionRec struct {
	ID      int
	TimeAgo string
}

type row struct {
	User    string
	Results []resultRec
}

type resultsModel struct {
	Sessions []sessionRec
	Rows     []row
	// Users    []string
	// Results  []resultRec //len of array = len(Sessions)*len(Users). string by string.
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

func NewResultsModel(DB *storm.DB) *resultsModel {
	res := &resultsModel{}
	now := time.Now()

	//fill sessions
	allSessions := (&repositories.Sessions{DB}).FindAll()
	for _, curSession := range *allSessions {
		res.Sessions = append(res.Sessions, sessionRec{curSession.ID, timeAgoHumanString(now, curSession.CreatedTime)})
	}

	//fill users
	allRuns := (&repositories.Runs{DB}).FindAll()
	users := map[string]bool{}
	for _, curRun := range *allRuns {
		if users[curRun.DeviceOwnerName] == false {
			users[curRun.DeviceOwnerName] = true
		}
	}

	for curUser, _ := range users {
		curRow := row{}
		curRow.User = curUser
		for _, curSession := range res.Sessions {
			resultRuns := (&repositories.Runs{DB}).FindBySessionIDAndUserName(curSession.ID, curUser)
			curSessionUserResult := resultRec{}
			if len(*resultRuns) == 0 {
				curSessionUserResult.Status = ResultStatusNotRunned
			} else {
				curSessionUserResult.TriesCount = len(*resultRuns)
				curSessionUserResult.LastTryAgo = timeAgoHumanString(now, (*resultRuns)[0].StartTime)
				if (*resultRuns)[0].ExitCode == `0` {
					curSessionUserResult.Status = ResultStatusLastPassed
				} else {
					curSessionUserResult.Status = ResultStatusLastLastFailed
				}
			}
			curRow.Results = append(curRow.Results, curSessionUserResult)
		}
		res.Rows = append(res.Rows, curRow)
	}

	// //fill results
	// for _, curUser := range res.Users {
	// 	for _, curSession := range res.Sessions {
	// 		resultRuns := (&repositories.Runs{DB}).FindBySessionIDAndUserName(curSession.ID, curUser)
	// 		curSessionUserResult := resultRec{}
	// 		if len(*resultRuns) == 0 {
	// 			curSessionUserResult.Status = ResultStatusNotRunned
	// 		} else {
	// 			curSessionUserResult.TriesCount = len(*resultRuns)
	// 			curSessionUserResult.LastTryAgo = timeAgoHumanString(now, (*resultRuns)[0].StartTime)
	// 			if (*resultRuns)[0].ExitCode == `0` {
	// 				curSessionUserResult.Status = ResultStatusLastPassed
	// 			} else {
	// 				curSessionUserResult.Status = ResultStatusLastLastFailed
	// 			}
	// 		}
	// 		res.Results = append(res.Results, curSessionUserResult)
	// 	}
	// }

	// allTestpacks := (&repositories.Testpacks{DB}).FindAll()
	// for _, curTestpack := range *allTestpacks {
	// 	res.Tesptacks = append(res.Tesptacks, testpackRec{curTestpack.ID, timeAgoHumanString(now, curTestpack.UploadTime)})
	// }
	log.Println(res)

	return res
}
