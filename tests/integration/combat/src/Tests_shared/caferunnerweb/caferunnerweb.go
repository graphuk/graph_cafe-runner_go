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

func (t *CafeRunnerWeb) OpenPageSession() *PageSession {
	return NewOpenPageSession(t.bro)
}

func (t *CafeRunnerWeb) Cleanup() {
	if t.bro != nil {
		t.bro.Cancel()
	}
}