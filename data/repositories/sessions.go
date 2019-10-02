package repositories

import (
	//	"log"
	"time"

	"github.com/asdine/storm"

	"github.com/graph-uk/graph_cafe-runner_go/data/models"
)

type Sessions struct {
	Tx storm.Node
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Sessions) Create() *models.Session {
	session := &models.Session{
		CreatedTime: time.Now(),
	}

	check(t.Tx.Save(session))
	return session
}

func (t *Sessions) Find(id int) *models.Session {
	res := &models.Session{}
	check(t.Tx.One(`ID`, id, res))
	return res
}

func (t *Sessions) FindAll() *[]models.Session {
	res := &[]models.Session{}
	check(t.Tx.All(res))
	return res
}

func (t *Sessions) FindAllOrderIDDesc() *[]models.Session {
	res := &[]models.Session{}
	query := t.Tx.Select().OrderBy(`ID`).Reverse()
	check(query.Find(res))
	return res
}

func (t *Sessions) FindByTestpackID(testpackID int) *[]models.Session {
	res := &[]models.Session{}
	err := t.Tx.Find(`TestpackID`, testpackID, res)
	if err != nil { //ignore "not found" error. Return empty slice. All other errors are critical.
		if err.Error() != `not found` {
			check(err)
		}
	}
	return res
}

func (t *Sessions) FindLast() *models.Session {
	res := &[]models.Session{}
	check(t.Tx.All(res, storm.Limit(1), storm.Reverse()))
	return &(*res)[0]
}
