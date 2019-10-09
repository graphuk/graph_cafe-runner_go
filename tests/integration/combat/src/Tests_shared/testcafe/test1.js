import CafeWrapper from './cafeWrapper'

fixture `Example test with env variables`


test('Keros layout guest', async t => {
    await t.navigateTo(process.env.HOSTNAME);
    let cafe = await new CafeWrapper(t);
    await cafe.resizeWindow(720, 800);
    await cafe.typeByXpath('//input[@id="text"]',process.env.REQUEST);
    await cafe.clickByXpath('//button[@type="submit"]');
    await cafe.checkShownByXpath("//li[@class='serp-item']//*[text()='"+process.env.RESULT_SITE_URL+"']")
})