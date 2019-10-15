package caferunnerweb

import (
	"Tests_shared/browser"
)

type PageResults struct {
	Bro        *browser.Browser
	PartHeader *PartHeader
}

func NewPageResults(bro *browser.Browser) *PageResults {
	res := PageResults{bro, NewPartHeader(bro)}
	return &res
}

//2nd column is first result column. 1st column contains device name.
func (t *PageResults) CheckCellClassByDeviceNameAndColumn(deviceName, column, class string) {
	t.Bro.ClickByXpath(`(//div[@class='rTableCell' and text()='` + deviceName + `']/../div[` + column + `])[@class='` + class + `']`)
}
