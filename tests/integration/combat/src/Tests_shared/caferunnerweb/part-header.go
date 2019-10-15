package caferunnerweb

import (
	"Tests_shared/browser"
)

type PartHeader struct {
	Bro *browser.Browser
}

func NewPartHeader(bro *browser.Browser) *PartHeader {
	res := PartHeader{bro}
	return &res
}

func (t *PartHeader) ClickResults() *PageResults {
	t.Bro.ClickByXpath(`//a[@href='/results']`)
	return NewPageResults(t.Bro)
}
