package services

import (
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/data"
	"github.com/graph-uk/graph_cafe-runner_go/data/models"
	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"
)

type Session struct {
	Tx storm.Node
}

func (t *Session) Create(TestpackID int) *models.Session {
	var res *models.Session
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		testpack := (&repositories.Testpacks{tx}).Find(TestpackID) // check testpack exist

		if testpack.Status != models.TPStatusReadyForRunning {
			//return nil, errors.New(`Failed to create session for testpack ` + TestpackID + `, because it not ready for running. Status of testpack is ` + strconv.Itoa(int(testpack.Status)))
			panic(`Failed to create session for testpack ` + strconv.Itoa(TestpackID) + `, because it not ready for running. Status of testpack is ` + strconv.Itoa(int(testpack.Status)))
		}

		res = (&repositories.Sessions{tx}).Create(TestpackID)

		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
	return res
}
