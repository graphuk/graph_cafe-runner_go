package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/graph-uk/graph_cafe-runner_go/cmd/cafe-runner-server/config"
	"github.com/graph-uk/graph_cafe-runner_go/data"
	"github.com/graph-uk/graph_cafe-runner_go/data/models"
	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"
	"github.com/graph-uk/graph_cafe-runner_go/logic/utils"
)

const runPathTemplate = "_data/runs/%d"

type Run struct {
	Tx               storm.Node
	CafeRunnerConfig *config.Configuration
}

func (t *Run) Create(SessionID int, TestpackID int, DeviceOwnerName string) *models.Run {
	var res *models.Run

	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		(&repositories.Sessions{tx}).Find(SessionID)   // check session exist
		(&repositories.Testpacks{tx}).Find(TestpackID) // check testpack exist
		res = (&repositories.Runs{tx}).Create(SessionID, TestpackID, DeviceOwnerName)
		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
	return res
}

func (t *Run) npmInstall(RunID int) {
	runPath := fmt.Sprintf(runPathTemplate, RunID)

	t.updateStatus(RunID, models.RunStatusReadyForNPMInstall, models.RunStatusNPMInstallProgress)

	out, err := utils.ExecuteCommand(`npm`, []string{`install`}, runPath, os.Environ(), time.Minute*3)
	if err != nil {
		log.Println(`Failed to run "npm install" for run ` + strconv.Itoa(RunID))
		log.Println(err.Error())
		log.Println(string(out))
		unitOfWork := data.UnitOfWork{t.Tx}
		command := func(tx storm.Node) (interface{}, error) {
			run := (&repositories.Runs{t.Tx}).Find(RunID)
			if run.Status != models.RunStatusCopyTestpackInProgress {
				panic(`Status of run ` + strconv.Itoa(RunID) + `is not expected. Should be CopyInProgress, but actual is ` + strconv.Itoa(int(run.Status)))
			}
			run.ExitCode = `1` //because it should be non-zero for error
			run.StdOut = []byte(err.Error())
			run.Status = models.RunStatusCopyTestpackFailed
			(&repositories.Runs{tx}).Update(run)
			return nil, nil
		}
		unitOfWork.ExecuteCommand(command)
		panic(err)
	}
	t.updateStatus(RunID, models.RunStatusNPMInstallProgress, models.RunStatusReadyForCafeThread)

	//t.updateStatus(RunID, models.TPStatusInitProgress, models.TPStatusReadyForRunning)
}

//let's method "isMyHostname" always return true
//otherwise, in docker, testcafe failed with message
//ERROR The specified "<someIP>" hostname cannot be resolved to the current machine.
//that's due to isMyHostname run server at random port (not forwarded to container), and cannot connect to itself.
func (t *Run) hackEndpointUtils(RunID int) {
	runPath := fmt.Sprintf(runPathTemplate, RunID)
	contentBytes, err := ioutil.ReadFile(runPath + `/node_modules/endpoint-utils/index.js`)
	check(err)
	contentString := string(contentBytes)
	contentString = strings.Replace(contentString, `    return getFreePort()`, `    return true; return getFreePort()`, -1)
	check(ioutil.WriteFile(runPath+`/node_modules/endpoint-utils/index.js`, []byte(contentString), 644))
}

//let's idle page (at the end of session) redirect back to results.
func (t *Run) hackIdlePageForBackRedirect(RunID int, externalURL string) {
	runPath := fmt.Sprintf(runPathTemplate, RunID)
	contentBytes, err := ioutil.ReadFile(runPath + `/node_modules/testcafe/lib/client/browser/idle-page/index.html.mustache`)
	check(err)
	contentString := string(contentBytes)
	contentString = strings.Replace(contentString, `    new IdlePage('{{{statusUrl}}}', '{{{heartbeatUrl}}}', '{{{initScriptUrl}}}', { retryTestPages: {{{retryTestPages}}} });`,
		`    new IdlePage('{{{statusUrl}}}', '{{{heartbeatUrl}}}', '{{{initScriptUrl}}}', { retryTestPages: {{{retryTestPages}}} });
    function updateRunStatus(){
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function () {
            if (this.readyState != 4) return;
            if (this.status != 200) {
                window.location.href="`+externalURL+`/runs/`+strconv.Itoa(RunID)+`";
                //window.history.go(-2);
            }
        };
        xhr.open('GET', "{{{statusUrl}}}", true);
        xhr.send();
    }
    setInterval(updateRunStatus, 5000);`, -1)
	check(ioutil.WriteFile(runPath+`/node_modules/testcafe/lib/client/browser/idle-page/index.html.mustache`, []byte(contentString), 644))
}

func (t *Run) copyTestpack(RunID int) {
	t.updateStatus(RunID, models.RunStatusReadyForCopyTestpack, models.RunStatusCopyTestpackInProgress)
	run := (&repositories.Runs{t.Tx}).Find(RunID)
	//session := (&repositories.Sessions{t.Tx}).Find(run.SessionID)
	runPath := fmt.Sprintf(runPathTemplate, RunID)

	err := (&Testpack{t.Tx}).CopyToFolder(run.TestpackID, runPath)
	if err != nil { // if copying failed
		t.markAsCopyTestpackFailed(RunID, err)
	} else { //if copying succeed
		t.updateStatus(RunID, models.RunStatusCopyTestpackInProgress, models.RunStatusReadyForNPMInstall)
		//t.npmInstall(RunID)
	}
}

func (t *Run) runCafeThread(RunID int) {
	freePort1, freePort2 := utils.GetFirstFreeLocalPorts(t.CafeRunnerConfig.Server.Cafe.LowPort, t.CafeRunnerConfig.Server.Cafe.HighPort)
	log.Println(`Run ` + strconv.Itoa(RunID) + `: free ports - ` + strconv.Itoa(freePort1) + `,` + strconv.Itoa(freePort2))
	runPath := fmt.Sprintf(runPathTemplate, RunID)

	pwd, err := os.Getwd()
	check(err)

	run := (&repositories.Runs{t.Tx}).Find(RunID)
	testpack := (&repositories.Testpacks{t.Tx}).Find(run.TestpackID)
	envVars := append(testpack.EnvVars, os.Environ()...)

	cmd, err := utils.StartCmd(`node`, []string{pwd + `/` + runPath + `/node_modules/testcafe/bin/testcafe.js`, `remote`, `test1.js`, `--hostname`, t.CafeRunnerConfig.Server.Hostname, `--ports`, strconv.Itoa(freePort1) + `,` + strconv.Itoa(freePort2)}, pwd+`/`+runPath, envVars)

	if err != nil {
		log.Println(`Failed to start cafe thread` + strconv.Itoa(RunID))
		log.Println(err.Error())
		log.Println(string(cmd.StdOutBuf))
		t.markAsCafeThreadFailed(RunID, err)
		panic(err)
	}

	t.markAsCafeThreadProgress(RunID, freePort1)
	go t.watchCafeThread(RunID, cmd, time.Minute*10)
}

func (t *Run) RunInitSteps(RunID int, ExternalURL string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR: Failed on init steps for run: ", RunID)
			log.Println(r)
			debug.PrintStack()
		}
	}()
	log.Println(`Run ` + strconv.Itoa(RunID) + `. Init started.`)
	t.copyTestpack(RunID)
	t.npmInstall(RunID)
	t.hackEndpointUtils(RunID)
	t.hackIdlePageForBackRedirect(RunID, ExternalURL)
	t.runCafeThread(RunID)
	log.Println(`Run ` + strconv.Itoa(RunID) + `. Init finished. Connect for testing.`)
}

func (t *Run) updateStatus(RunID int, statusFrom, statusTo uint8) {
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		run := (&repositories.Runs{t.Tx}).Find(RunID)
		if run.Status != statusFrom {
			panic(`ERROR: tespack ` + strconv.Itoa(RunID) + ` has status ` + strconv.Itoa(int(run.Status)) + ` but it should be ` + strconv.Itoa(int(statusFrom)))
		}

		run.Status = statusTo
		(&repositories.Runs{tx}).Update(run)

		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
}

func (t *Run) markAsCopyTestpackFailed(RunID int, err error) {
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		run := (&repositories.Runs{t.Tx}).Find(RunID)
		if run.Status != models.RunStatusCopyTestpackInProgress {
			panic(`Status of run ` + strconv.Itoa(RunID) + `is not expected. Should be CopyInProgress, but actual is ` + strconv.Itoa(int(run.Status)))
		}
		run.ExitCode = `1` //because it should be non-zero for error
		run.StdOut = []byte(err.Error())
		run.Status = models.RunStatusCopyTestpackFailed
		(&repositories.Runs{tx}).Update(run)
		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
}

func (t *Run) markAsCafeThreadFailed(RunID int, err error) {
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		run := (&repositories.Runs{t.Tx}).Find(RunID)
		if run.Status != models.RunStatusReadyForCafeThread {
			panic(`Status of run ` + strconv.Itoa(RunID) + `is not expected. Should be ReadyForCafeThread, but actual is ` + strconv.Itoa(int(run.Status)))
		}
		run.ExitCode = `1` //because it should be non-zero for error
		run.StdOut = []byte(err.Error())
		run.Status = models.RunStatusCafeThreadFailed

		(&repositories.Runs{tx}).Update(run)
		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
}

func (t *Run) markAsCafeThreadProgress(RunID, port int) {
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		run := (&repositories.Runs{t.Tx}).Find(RunID)
		if run.Status != models.RunStatusReadyForCafeThread {
			panic(`Status of run ` + strconv.Itoa(RunID) + `is not expected. Actual is ` + strconv.Itoa(int(run.Status)))
		}
		run.Port = port
		run.Status = models.RunStatusCafeThreadProgress

		(&repositories.Runs{tx}).Update(run)
		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
}

func (t *Run) markAsCafeThreadTimeout(RunID int) {
	t.updateStatus(RunID, models.RunStatusCafeThreadProgress, models.RunStatusCafeThreadFinTimeout)
}

func (t *Run) markAsCafeThreadFin(RunID int, exitCode int, stdout *[]byte) {
	unitOfWork := data.UnitOfWork{t.Tx}
	command := func(tx storm.Node) (interface{}, error) {
		run := (&repositories.Runs{t.Tx}).Find(RunID)
		if run.Status != models.RunStatusCafeThreadProgress {
			panic(`Status of run ` + strconv.Itoa(RunID) + `is not expected. Actual is ` + strconv.Itoa(int(run.Status)))
		}
		run.ExitCode = strconv.Itoa(exitCode) //because it should be non-zero for error
		run.StdOut = *stdout
		run.Status = models.RunStatusCafeThreadFin

		(&repositories.Runs{tx}).Update(run)
		return nil, nil
	}
	unitOfWork.ExecuteCommand(command)
}

func (t *Run) watchCafeThread(RunID int, cmd *utils.CmdProcess, timeout time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			//debug.PrintStack()
		}
	}()
	exitCode, err := cmd.WaitingForExitCode(timeout)
	if err != nil {
		t.markAsCafeThreadTimeout(RunID)
		cmd.KillWithChilds()
		//log.Println(string(append(cmd.StdOutBuf, cmd.StdErrBuf...)))
		panic(`Run ` + strconv.Itoa(RunID) + `. Cafe thread timeout was reached.`)
	}
	combineOut := append(cmd.StdOutBuf, cmd.StdErrBuf...)
	t.markAsCafeThreadFin(RunID, exitCode, &combineOut)
	if exitCode == 0 {
		log.Println(`Run ` + strconv.Itoa(RunID) + `. Cafe thread finished with exitCode 0`)
	} else {
		log.Println(`Run ` + strconv.Itoa(RunID) + `. Cafe thread finished with exitCode ` + strconv.Itoa(exitCode))
		log.Println(`Run ` + strconv.Itoa(RunID) + `. Cafe thread StdOut: ` + string(combineOut))
	}
}
