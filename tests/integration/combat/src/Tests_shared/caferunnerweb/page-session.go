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

// //2nd column is first result column. 1st column contains device name.
// func (t *PageSession) CheckCellClassByDeviceNameAndColumn(deviceName, column, class string) {
// 	t.Bro.ClickByXpath(`(//div[@class='rTableCell' and text()='` + deviceName + `']/../div[` + column + `])[@class='` + class + `']`)
// }
