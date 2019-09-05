package services

import (
	"fmt"
	"io/ioutil"
	"log"

	"runtime/debug"
	"strconv"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/data"
	"github.com/graph-uk/graph_cafe-runner_go/data/models"
	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"
	"github.com/graph-uk/graph_cafe-runner_go/logic/utils"
)

const testpackArchivePathTemplate = "_data/testpacks/%d.zip"
const testpackUnarchivedPathTemplate = "_data/testpacks/%d"

type Testpack struct {
	Tx storm.Node
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (t *Testpack) RunInitSteps(TestpackID int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed on init steps for tesptack: ", TestpackID)
			debug.PrintStack()
		}
	}()
	log.Println(`Testpack ` + strconv.Itoa(TestpackID) + `. Init started.`)
	t.unzip(TestpackID)
	//t.init(TestpackID)
	log.Println(`Testpack ` + strconv.Itoa(TestpackID) + `. Init finished.`)
}

func (t *Testpack) unzip(TestpackID int) {
	t.updateStatus(TestpackID, models.TPStatusReadyForUnzip, models.TPStatusUnzipInProgress)

	//save zip from DB to FS
	testpack := (&repositories.Testpacks{t.Tx}).Find(TestpackID)
	testpackPathArchive := fmt.Sprintf(testpackArchivePathTemplate, testpack.ID)
	check(ioutil.WriteFile(testpackPathArchive, testpack.Zip, 0666))

	//unzip zipfile on FS
	testpackPathUnarchived := fmt.Sprintf(testpackUnarchivedPathTemplate, testpack.ID)
	check(utils.Unzip(testpackPathArchive, testpackPathUnarchived))

	t.updateStatus(TestpackID, models.TPStatusUnzipInProgress, models.TPStatusReadyForRunning)
}

// func (t *Testpack) init(TestpackID int) {
// 	t.updateStatus(TestpackID, models.TPStatusReadyForInit, models.TPStatusInitProgress)

// 	// if false {
// 	// 	testpackPathUnarchived := fmt.Sprintf(testpackUnarchivedPathTemplate, TestpackID)
// 	// 	out, err := utils.ExecuteCommand(`npm`, []string{`install`}, testpackPathUnarchived+`/testcafe`, os.Environ(), time.Minute)
// 	// 	if err != nil {
// 	// 		log.Println(`Failed to init testpack ` + strconv.Itoa(TestpackID))
// 	// 		log.Println(err.Error())
// 	// 		log.Println(string(out))
// 	// 		t.markAsInitFailed(TestpackID, out)
// 	// 		panic(err)
// 	// 	}
// 	// }

// 	t.updateStatus(TestpackID, models.TPStatusInitProgress, models.TPStatusReadyForRunning)
// }

func (t *Testpack) updateStatus(TestpackID int, statusFrom, statusTo uint8) {
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		testpack := (&repositories.Testpacks{t.Tx}).Find(TestpackID)
		if testpack.Status != statusFrom {
			panic(`ERROR: tespack ` + strconv.Itoa(TestpackID) + ` has status ` + strconv.Itoa(int(testpack.Status)) + ` but it should be ` + strconv.Itoa(int(statusFrom)))
		}

		testpack.Status = statusTo
		(&repositories.Testpacks{tx}).Update(testpack)

		// if statusTo == models.TPStatusReadyForRunning { // test nested transactions
		// 	(&Session{tx}).Create(TestpackID)
		// }

		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
}

// func (t *Testpack) markAsInitFailed(TestpackID int, failOut []byte) {
// 	unitOfWork := data.UnitOfWork{t.Tx}
// 	command := func(tx storm.Node) (interface{}, error) {
// 		testpack := (&repositories.Testpacks{t.Tx}).Find(TestpackID)

// 		testpack.InitFailOut = failOut
// 		testpack.Status = models.TPStatusInitFailed

// 		(&repositories.Testpacks{tx}).Update(testpack)
// 		return nil, nil
// 	}
// 	unitOfWork.ExecuteCommand(command)
// }

func (t *Testpack) CopyToFolder(targetFolder string) error {
	testpackPathUnarchived := fmt.Sprintf(testpackUnarchivedPathTemplate)
	return utils.CopyDir(testpackPathUnarchived, targetFolder)
}
