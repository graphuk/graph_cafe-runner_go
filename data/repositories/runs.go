package repositories

import (
	"strconv"
	"time"

	"github.com/asdine/storm/q"

	"github.com/asdine/storm"

	"github.com/graph-uk/graph_cafe-runner_go/data/models"
)

type Runs struct {
	Tx storm.Node
}

func (t *Runs) Create(SessionID int, TestpackID int, DeviceOwnerName string) *models.Run {
	session := (&Sessions{t.Tx}).Find(SessionID)
	if session == nil {
		panic(`Session not found by ID: ` + strconv.Itoa(SessionID))
	}

	testpack := (&Testpacks{t.Tx}).Find(TestpackID)
	if testpack == nil {
		panic(`Testpack not found by ID: ` + strconv.Itoa(TestpackID))
	}

	run := &models.Run{
		Status:          models.RunStatusReadyForCopyTestpack,
		SessionID:       SessionID,
		TestpackID:      TestpackID,
		StartTime:       time.Now(),
		DeviceOwnerName: DeviceOwnerName,
	}

	check(t.Tx.Save(run))

	return run
}

func (t *Runs) Find(id int) *models.Run {
	res := &models.Run{}
	check(t.Tx.One(`ID`, id, res))
	return res
}

func (t *Runs) FindAll() *[]models.Run {
	res := &[]models.Run{}
	check(t.Tx.All(res))
	return res
}

func (t *Runs) FindBySessionID(sessionID int) *[]models.Run {
	res := &[]models.Run{}
	err := t.Tx.Find(`SessionID`, sessionID, res)
	if err != nil { //ignore "not found" error. Return empty slice. All other errors are critical.
		if err.Error() != `not found` {
			check(err)
		}
	}
	return res
}

func (t *Runs) FindBySessionIDAndUserName(sessionID int, userName string) *[]models.Run {
	res := &[]models.Run{}
	query := t.Tx.Select(q.And(q.Eq(`SessionID`, sessionID), q.Eq(`DeviceOwnerName`, userName)))
	err := query.Find(res)
	if err != nil { //ignore "not found" error. Return empty slice. All other errors are critical.
		if err.Error() != `not found` {
			check(err)
		}
	}
	return res
}

func (t *Runs) Update(run *models.Run) {
	check(t.Tx.Update(run))
}
