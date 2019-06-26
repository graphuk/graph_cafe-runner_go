package browser

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"

	"github.com/chromedp/chromedp"
)

type Browser struct {
	ctx     context.Context
	Cancel  context.CancelFunc
	Timeout time.Duration
	Retries int
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func timeout(f func(), timeout time.Duration) {
	ch := make(chan bool, 2)

	timeoutReached := new(bool)
	*timeoutReached = false

	panicInFunc := new(bool)
	*panicInFunc = false

	go func() {
		defer func() { //if
			if r := recover(); r != nil {
				*panicInFunc = true
			}
			ch <- true
			close(ch)
		}()

		f()

	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-ch:
		if *panicInFunc {
			panic(`Panic in func before timeout`)
		}
	case <-timer.C:
		*timeoutReached = true
		panic(`Timeout was reached.`)
	}
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func NewBrowser() *Browser {
	res := Browser{}

	execAllocator := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		//chromedp.Headless,

		// After Puppeteer's default behavior.
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		//chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),
		//chromedp.Flag("user-data-dir", dataDir),
		//chromedp.Flag("load-extension", `"`+dataDir+`\uVPN - free and unlimited Chrome VPN.crx"`), //strange way to chrome binary
		//chromedp.Flag("disable-extensions-file-access-check", true),
		//chromedp.Flag("allow-legacy-extension-manifests", true),
		//--disable-extensions-file-access-check --
		//chromedp.Flag("load-extension", `vpn`),
	}

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), execAllocator...)

	// also set up a custom logger
	res.ctx, res.Cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	res.Timeout = 120 * time.Second
	res.Retries = 3

	return &res
}

func (t *Browser) ClickByXpath(xpath string) {
	log.Println(`ClickByXpath ` + xpath)

	timeout(func() {
		check(chromedp.Run(t.ctx, chromedp.Click(xpath, chromedp.NodeVisible)))
	}, t.Timeout)
}

func (t *Browser) FillByXpath(xpath, value string) {
	log.Println(`ClickByXpath ` + xpath)

	timeout(func() {
		check(chromedp.Run(t.ctx, chromedp.SendKeys(xpath, value, chromedp.NodeVisible)))
	}, t.Timeout)
}

func (t *Browser) Get(url string) {
	log.Println(`Get ` + url)

	timeout(func() {
		check(chromedp.Run(t.ctx, chromedp.Navigate(url)))
	}, t.Timeout)
}

func (t *Browser) Refresh() {
	log.Println(`Refresh page`)

	timeout(func() {
		check(chromedp.Run(t.ctx, chromedp.Reload()))
	}, t.Timeout)
}

func (t *Browser) GetTextByXpath(xpath string) *[]string {
	log.Println(`GetTextByXpath ` + xpath)
	//panic(`fakePanic`)

	res := []string{}

	timeout(func() {
		var nodes []*cdp.Node

		check(chromedp.Run(t.ctx,
			chromedp.Nodes(xpath, &nodes),
		))

		for _, curNode := range nodes {
			attr := ``
			check(chromedp.Run(t.ctx, chromedp.Text(curNode.FullXPath(), &attr)))
			res = append(res, attr)
		}
	}, t.Timeout)
	return &res
}

func (t *Browser) GetAttributeValueByXpath(xpath, attr string) *[]string {
	log.Println(`GetAttributeValueByXpath ` + xpath + ` attr ` + attr)
	res := []string{}

	timeout(func() {
		var nodes []*cdp.Node

		check(chromedp.Run(t.ctx,
			chromedp.Nodes(xpath, &nodes),
		))

		for _, curNode := range nodes {
			res = append(res, curNode.AttributeValue(attr))
		}
	}, t.Timeout)
	return &res
}

func (t *Browser) GetNodesByXpath(xpath string) []*cdp.Node {
	log.Println(`GetNodesByXpath ` + xpath)
	var nodes []*cdp.Node
	timeout(func() {

		check(chromedp.Run(t.ctx,
			chromedp.Nodes(xpath, &nodes),
		))
	}, t.Timeout)
	return nodes
}
