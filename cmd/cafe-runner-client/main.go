package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/graph-uk/graph_cafe-runner_go/logic/utils"
)

func compressTestpack() (string, error) {
	fmt.Println("Compressing testpack")
	tmpFile, err := ioutil.TempFile("", "cafeSession")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()
	curpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	utils.Zipit(curpath, tmpFile.Name())
	return tmpFile.Name(), nil
}

func main() {

	zipName, _ := compressTestpack()

	zipbytes, err := ioutil.ReadFile(zipName)
	if err != nil {
		panic(err)
	}

	content := base64.StdEncoding.EncodeToString(zipbytes)
	json := fmt.Sprintf("{\"Content\": \"%s\"}", content)
	body := bytes.NewBuffer([]byte(json))

	serverAddress := os.Args[1]
	resp, err := http.Post(serverAddress+`/api/v1/testpacks`, "application/json", body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
}
