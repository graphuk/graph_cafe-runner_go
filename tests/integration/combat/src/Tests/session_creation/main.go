package main

import (
	"Tests_shared/aTest"
	"Tests_shared/caferunnerweb"

	"Tests_shared/cmdutils"

	"fmt"
	"os"
	"time"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		HostName aTest.StringParam
	}
	timestamp string
}

func createNewTest() *theTest {
	var result theTest
	result.aTest.DefaultParams = append(result.aTest.DefaultParams, "-HostName=localhost:3133")
	result.aTest.Tags = append(result.aTest.Tags, "local")
	result.aTest.FillParamsFromCLI(&result.params)
	result.timestamp = time.Now().Format("20060102150405")
	fmt.Println("Timestamp: " + result.timestamp)
	result.aTest.CreateOutputFolder()
	return &result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	theTest := createNewTest()
	fmt.Println(theTest.timestamp)

	defer func() {
		if r := recover(); r != nil {
			aTest.PrintSourceAndContinuePanic(r)
		}
	}()

	pwd, err := os.Getwd()
	check(err)

	check(cmdutils.CopyFile(`../../Tests_shared/cafe-runner-server.exe`, `out/cafe-runner-server.exe`))

	server := cmdutils.MustStartCmd(pwd+`/out/cafe-runner-server.exe`, []string{}, pwd+`/out`, os.Environ())
	defer server.KillWithChilds()
	server.WaitingForStdOutContains(`http server started on`, time.Second*10)
	cmdutils.DownloadFile(pwd+`/out/cafe-runner-client.exe`, `http://`+theTest.params.HostName.Value+`/assets/assets/dist/win64/cafe-runner-client.exe`)

	client := cmdutils.MustStartCmd(pwd+`/out/cafe-runner-client.exe`, []string{`http://127.0.0.1:3133`}, pwd+`/../../Tests_shared/testcafe`, os.Environ())

	server.WaitingForStdErrContains(`Post testpack received.`, time.Second*30)
	server.WaitingForStdErrContains(`Testpack created with id: 1`, time.Second*10)
	client.WaitingForExitWithCode(0, time.Second*10)

	server.WaitingForStdErrContains(`Testpack 1. Init started.`, time.Second*10)
	server.WaitingForStdErrContains(`Testpack 1. Init finished.`, time.Second*10)

	cmdutils.MustPost(`http://`+theTest.params.HostName.Value+`/api/v1/sessions`, `{}`)

	server.WaitingForStdErrContains(`Post session received.`, time.Second*10)
	server.WaitingForStdErrContains(`Session is created with id: 1`, time.Second*10)

	fmt.Println(string(server.StdErrBuf))

	return

	cafeRunnerWeb := caferunnerweb.NewCafeRunnerWeb()
	defer cafeRunnerWeb.Cleanup()
	PageSession := cafeRunnerWeb.OpenPageSession()
	PageSession.FillDeviceOwnerName(`test1`)
	PageSession.ClickStartTesting()

	server.WaitingForStdErrContains(`Post run received.`, time.Second*10)
	server.WaitingForStdErrContains(`Run for session 1 created with id: 1`, time.Second*10)
	server.WaitingForStdErrContains(`Run 1. Init started.`, time.Second*10)
	server.WaitingForStdErrContains(`Get run received.`, time.Second*10)
	server.WaitingForStdErrContains(`Run 1: free ports -`, time.Second*70)
	server.WaitingForStdErrContains(`remote test1.js --hostname localhost --ports`, time.Second*10)
	server.WaitingForStdErrContains(`Run 1. Init finished. Connect for testing.`, time.Second*10)
	server.WaitingForStdErrContains(`Run 1. Cafe thread finished with exitCode 0`, time.Second*60)
}
