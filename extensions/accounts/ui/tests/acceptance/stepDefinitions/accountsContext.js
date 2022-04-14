const assert = require('assert')
const { client } = require('nightwatch-api')
const { Given, When, Then } = require('@cucumber/cucumber')

When('the user browses to the accounts page', function () {
  return client.page.accountsPage().navigateAndWaitUntilMounted()
})

Then('user {string} should be displayed in the accounts list on the WebUI', async function (username) {
  await client.page.accountsPage().accountsList()
  const userListed = await client.page.accountsPage().isUserListed(username)
  return assert.strictEqual(userListed, true)
})

Then('user {string} should not be displayed in the accounts list on the WebUI', async function (username) {
  await client.page.accountsPage().accountsList()
  const userDeleted = await client.page.accountsPage().isUserDeleted(username)
  return assert.strictEqual(userDeleted, true)
})

Given('the user has changed the role of user {string} to {string}', function (username, role) {
  return client.page.accountsPage().selectRole(username, role)
})

When('the user changes the role of user {string} to {string} using the WebUI', function (username, role) {
  return client.page.accountsPage().selectRole(username, role)
})

Then('the displayed role of user {string} should be {string} on the WebUI', function (username, role) {
  return client.page.accountsPage().checkUsersRole(username, role)
})

Then('the user should not be able to see the accounts list on the WebUI', async function () {
  return client.page.accountsPage()
    .waitForAjaxCallsToStartAndFinish()
    .waitForElementVisible('@loadingAccountsListFailed')
})

When('the user disables user/users {string} using the WebUI', function (usernames) {
  return client.page.accountsPage().setUserActivated(usernames, false)
})

When('the user enables user/users {string} using the WebUI', function (usernames) {
  return client.page.accountsPage().setUserActivated(usernames, true)
})

Then('the status indicator of user/users {string} should be {string} on the WebUI', function (usernames, status) {
  return client.page.accountsPage().checkUsersStatus(usernames, status)
})

When(
  'the user creates a new user with username {string}, email {string} and password {string} using the WebUI',
  function (username, email, password) {
    return client.page.accountsPage().createUser(username, email, password)
  }
)

When('the user deletes user/users {string} using the WebUI', function (usernames) {
  return client.page.accountsPage().deleteUsers(usernames)
})
