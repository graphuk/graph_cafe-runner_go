import CafeWrapper from './cafeWrapper'

fixture `EmptyFixture`.page(process.env.TEST_HOSTNAME);

test('Example test with env variables', async t => {
    let cafe = await new CafeWrapper(t);
    await cafe.typeByXpath('//input[@name="text"]',process.env.TEST_REQUEST);
    await cafe.clickByXpath('//button[@type="submit"]');
    await cafe.checkShownByXpath("//*[@class='serp-item']//*[text()='"+process.env.TEST_RESULT_SITE_URL+"SomethingThatNotExistOnThePage']")
})