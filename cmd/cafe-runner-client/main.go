package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/graph-uk/graph_cafe-runner_go/api/testpacks/models"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//env variables may be provided after server address
//separator is space. format
//VAR1=VALUE1 VAR2=VALUE2 ...
func parseEnvVariables() []string {
	res := []string{}
	countOfArgs := len(os.Args)
	if countOfArgs > 2 {
		for i := 2; i < countOfArgs; i++ {
			res = append(res, os.Args[i])
		}
	}
	return res
}

func main() {
	zipName, err := compressTestpack()
	check(err)

	zipbytes, err := ioutil.ReadFile(zipName)
	check(err)

	content := base64.StdEncoding.EncodeToString(zipbytes)

	model := testpacks.TestpackPostModel{content, parseEnvVariables()}
	jsonBytes, err := json.Marshal(model)
	check(err)
	body := bytes.NewBuffer(jsonBytes)

	serverAddress := os.Args[1]

	resp, err := http.Post(serverAddress+`/api/v1/testpacks`, "application/json", body)
	check(err)
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	check(err)
}
