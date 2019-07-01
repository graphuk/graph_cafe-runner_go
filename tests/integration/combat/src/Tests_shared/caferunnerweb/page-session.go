package caferunnerweb

import (
	"Tests_shared/browser"
)

type PageSession struct {
	Bro *browser.Browser
}

func NewOpenPageSession(bro *browser.Browser) *PageSession {
	res := PageSession{bro}
	res.Bro.Get(`http://localhost:3133`)
	return &res
}

func (t *PageSession) FillDeviceOwnerName(value string) {
	t.Bro.FillByXpath(`//input`, value)
}

func (t *PageSession) ClickStartTesting() {
	t.Bro.ClickByXpath(`//button`)
}
