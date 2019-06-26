package repositories

import (
	"strconv"
	"time"

	"github.com/asdine/storm"

	"github.com/graph-uk/graph_cafe-runner_go/data/models"
)

type Runs struct {
	Tx storm.Node
}

func (t *Runs) Create(SessionID int, DeviceOwnerName string) *models.Run {
	session := (&Sessions{t.Tx}).Find(SessionID)

	if session == nil {
		panic(`Session not found by ID: ` + strconv.Itoa(SessionID))
	}

	run := &models.Run{
		Status:          models.RunStatusReadyForCopyTestpack,
		SessionID:       SessionID,
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

func (t *Runs) Update(run *models.Run) {
	check(t.Tx.Update(run))
}
