package main

import (
	"github.com/graph-uk/graph_cafe-runner_go/cmd/cafe-runner-server/config"
)

func main() {
	cafeRunnerConfig := config.GetConfiguration()
	cafeRunnerConfig.PrintToLog()

	cafeRunnerServer := &CafeRunnerServer{cafeRunnerConfig}
	cafeRunnerServer.Start()
}
