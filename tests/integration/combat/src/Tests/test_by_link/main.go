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

	client := cmdutils.MustStartCmd(pwd+`/out/cafe-runner-client.exe`, []string{`http://127.0.0.1:3133`, `HOSTNAME=http://ya.ru`, `REQUEST=github`, `RESULT_SITE_URL=github.com`}, pwd+`/../../Tests_shared/testcafe-success`, os.Environ())
	client.WaitingForExitWithCode(0, time.Second*10)

	server.WaitingForStdErrContains(`Testpack 1. Init finished.`, time.Second*90)

	cmdutils.MustPost(`http://`+theTest.params.HostName.Value+`/api/v1/sessions`, `{}`)

	server.WaitingForStdErrContains(`Session is created with id: 1`, time.Second*10)

	cafeRunnerWeb := caferunnerweb.NewCafeRunnerWeb()
	defer cafeRunnerWeb.Cleanup()
	PageRuntests := cafeRunnerWeb.OpenPageRuntests()
	PageRuntests.FillDeviceOwnerName(`test1`)

	link := PageRuntests.GetTestLink()

	PageRuntests.Bro.Get(link)

	server.WaitingForStdErrContains(`Run 1. Cafe thread finished with exitCode 0`, time.Second*90)
}
