const assert = require('assert')
const path = require('path')
const fs = require('fs-extra')
const { client } = require('nightwatch-api')
const { Given, When, Then, After } = require('cucumber')
const languageHelper = require('../helpers/language')

Given('the user browses to the settings page', function () {
  return client.page.settingsPage().navigateAndWaitTillLoaded()
})

Then('the setting {string} should have value {string}', async function (setting, result) {
  const actual = await client.page.settingsPage().getSettingsValue(setting)
  assert.strictEqual(actual, result, 'The setting value doesnt matches to ' + result)
})

When('the user changes the language to {string}', async function (value) {
  await client.page.settingsPage().changeSettings('Language', value)
})

Then('the files menu should be listed in language {string}', async function (language) {
  const menu = await client.page.filesPageSettingsContext().getMenuList()
  const expected = languageHelper.getFilesMenuForLanguage(language)
  assert.deepEqual(menu, expected, 'the menu list were not same')
})

Then('the account menu should be listed in language {string}', async function (language) {
  const menu = await client.page.filesPageSettingsContext().getUserMenu()
  const expected = languageHelper.getUserMenuForLanguage(language)
  assert.deepEqual(menu, expected, 'the menu list were not same')
})

Then('the files header should be displayed in language {string}', async function (language) {
  const items = await client.page.filesPageSettingsContext().getFileHeaderItems()
  const expected = languageHelper.getFilesHeaderMenuForLanguage(language)
  assert.deepEqual(items, expected, 'the menu list were not same')
})

After(async function () {
  const directory = path.join(client.globals.settings_store, 'values')
  try {
    console.log('Elements')
    fs.readdirSync(directory).map(element => {
      console.log(element)
    })
  } catch (err) {
    console.log('Error while reading the settings values from file system... ')
  }
  try {
    fs.emptyDirSync(directory)
  } catch (err) {
    console.log('Error while clearing settings values from file system')
    console.log('No settings may have been changed by the tests')
  }
})
