package caferunnerweb

import (
	"Tests_shared/browser"
)

type PageRuntests struct {
	Bro        *browser.Browser
	PartHeader *PartHeader
}

func NewPageRuntests(bro *browser.Browser) *PageRuntests {
	res := PageRuntests{bro, NewPartHeader(bro)}
	return &res
}

func OpenPageRuntests(bro *browser.Browser) {
	bro.Get(`http://localhost:3133`)
}

func (t *PageRuntests) FillDeviceOwnerName(value string) {
	t.Bro.FillByXpath(`//input`, value)
}

func (t *PageRuntests) ClickStartTesting() {
	t.Bro.ClickByXpath(`//button`)
}
