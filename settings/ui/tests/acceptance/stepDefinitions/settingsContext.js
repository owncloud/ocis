const assert = require('assert')
const path = require('path')
const fs = require('fs-extra')
const { client } = require('nightwatch-api')
const { Given, When, Then, After, Before } = require('cucumber')
const languageHelper = require('../helpers/language')

const initialLanguageAssignments = []

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
  console.log('def:acc')
  const menu = await client.page.filesPageSettingsContext().getUserMenu()
  console.log(menu)
  const expected = languageHelper.getUserMenuForLanguage(language)
  assert.deepStrictEqual(menu, expected, 'the menu list were not same')
})

Then('the files header should be displayed in language {string}', async function (language) {
  const items = await client.page.filesPageSettingsContext().getFileHeaderItems()
  const expected = languageHelper.getFilesHeaderMenuForLanguage(language)
  assert.deepStrictEqual(items, expected, 'the menu list were not same')
})

After(async function () {
  let directory = path.join(client.globals.settings_store, 'assignments')
  try {
    fs.readdirSync(directory).map(element => {
      if (!initialLanguageAssignments.includes(element)) {
        fs.unlinkSync(path.join(client.globals.settings_store, 'assignments', element))
      }
    })
  } catch (err) {
    console.log('Error while reading the settings values from file system... ')
  }

  directory = path.join(client.globals.settings_store, 'values')
  try {
    fs.emptyDirSync(directory)
  } catch (err) {
    console.log('Error while cleaning the settings values from file system... ')
  }
})

Before(async function() {
  const directory = path.join(client.globals.settings_store, 'assignments')
  try {
    fs.readdirSync(directory).map(element => {
      initialLanguageAssignments.push(element)
    })
  } catch (err) {
    console.log('Error while reading the settings values from file system... ')
  }
})