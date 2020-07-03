const assert = require('assert')
const { client } = require('nightwatch-api')
const { When, Then } = require('cucumber')

When('the user browses to the accounts page', function () {
  return client.page.accountsPage().navigateAndWaitTillLoaded()
})

Then('user {string} should be displayed in the accounts list on the WebUI', async function (username) {
  await client.page.accountsPage().accountsList(username)
  const userListed = await client.page.accountsPage().isUserListed(username)
  return assert.strictEqual(userListed, username)
})
