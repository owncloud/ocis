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

When('the user changes the role of user {string} to {string} using the WebUI', function (username, role) {
  return client.page.accountsPage().selectRole(username, role)
})

Then('the displayed role of user {string} should be {string} on the WebUI', function (username, role) {
  return client.page.accountsPage().checkUsersRole(username, role)
})

When('the user reloads the current page of the webUI', function () {
  return client.refresh()
})
