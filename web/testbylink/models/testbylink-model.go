package models

//	"github.com/graph-uk/graph_cafe-runner_go/data/repositories"

type testbylinkModel struct {
	SessionID  string
	TestpackID string
	Username   string
}

func NewRuntestsModel(sessionID, testpackID, username string) *testbylinkModel {
	return &testbylinkModel{sessionID, testpackID, username}
}
