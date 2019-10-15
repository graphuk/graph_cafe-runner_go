package caferunnerweb

import (
	"Tests_shared/browser"
)

type CafeRunnerWeb struct {
	bro *browser.Browser
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func NewCafeRunnerWeb() *CafeRunnerWeb {
	return &CafeRunnerWeb{browser.NewBrowser()}
}

func (t *CafeRunnerWeb) OpenPageRuntests() *PageRuntests {
	OpenPageRuntests(t.bro)
	return NewPageRuntests(t.bro)
}

func (t *CafeRunnerWeb) Cleanup() {
	if t.bro != nil {
		t.bro.Cancel()
	}
}
