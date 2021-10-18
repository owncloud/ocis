const assert = require('assert')
const { client } = require('nightwatch-api')
const { Given, When, Then } = require('@cucumber/cucumber')
const languageHelper = require('../helpers/language')

Given('the user browses to the settings page', function () {
  return client.page.settingsPage().navigateAndWaitTillLoaded()
})

Then('the setting {string} should have value {string}', async function (setting, result) {
  const actual = await client.page.settingsPage().getSettingsValue(setting)
  assert.strictEqual(actual, result, 'The setting value doesnt matches to ' + result)
})

Then('the setting {string} should not have any value', async function (setting) {
  const actual = await client.page.settingsPage().getSettingsValue(setting)
  assert.strictEqual(actual, false, 'The setting value was expected not to be present but was')
})

When('the user changes the language to {string}', async function (value) {
  await client.page.settingsPage().changeSettings('Language', value)
})

Then('the files menu should be listed in language {string}', async function (language) {
  const menu = await client.page.filesPageSettingsContext().getMenuList()
  const expected = languageHelper.getFilesMenuForLanguage(language)
  assert.deepStrictEqual(menu, expected, 'the menu list were not same')
})

Then('the account menu should be listed in language {string}', async function (language) {
  const menu = await client.page.filesPageSettingsContext().getUserMenu()
  const expected = languageHelper.getUserMenuForLanguage(language)
  assert.deepStrictEqual(menu, expected, 'the menu list were not same')
})

Then('the files header should be displayed in language {string}', async function (language) {
  const items = await client.page.filesPageSettingsContext().getFileHeaderItems()
  const expected = languageHelper.getFilesHeaderMenuForLanguage(language)
  assert.deepStrictEqual(items, expected, 'the menu list were not same')
})
