import { Selector } from 'testcafe';
import { ClientFunction } from 'testcafe';

const elementByXPath = Selector(xpath => {
    const iterator = document.evaluate(xpath, document, null, XPathResult.UNORDERED_NODE_ITERATOR_TYPE, null );
    const items = [];

    let item = iterator.iterateNext();

    while (item) {
        items.push(item);
        item = iterator.iterateNext();
    }

    return items;
});



export default class CafeWrapper {
    constructor (cafeController) {
        this.cafeController = cafeController;
    }

    async resizeWindow(width, height){
        await this.cafeController.resizeWindow(width, height);
    }

    async takeScreenshot(){
        await this.cafeController.takeScreenshot(new Date().getTime()+'.png');
    }

    async clickByXpath(xpath) {
        await this.takeScreenshot();
        await this.cafeController.click(elementByXPath(xpath));
    }

    async checkShownByXpath(xpath) {
        await elementByXPath(xpath).visible;
    }

    async typeByXpath(xpath, text) {
        await this.takeScreenshot();
        await this.cafeController.typeText(elementByXPath(xpath), text);
    }

    async goBack() {
        await this.takeScreenshot();
        const goBack = ClientFunction(() => window.history.back());
        await goBack();        
    }

    async getURL() {
        const getURL = ClientFunction(() => window.location.href);
        return await getURL();   
    }

    async checkUrlContains(urlPart){
        let url = await this.getURL();
        if (!url.includes(urlPart)){
            throw 'URL should contain "'+urlPart+'" substring, but not: ' + url;
        }
    }

    async getTextByXpath(xpath){
       return await elementByXPath(xpath).innerText;
    }

    async switchToIframeByXpath(xpath){
        await this.cafeController.switchToIframe(elementByXPath(xpath));
    }

    async switchToMainWindow(){
        await this.cafeController.switchToMainWindow();
    }

    async refreshPage(){
        await this.cafeController.eval(() => location.reload(true));
        //await this.cafeController.refreshPage();
    }

}
