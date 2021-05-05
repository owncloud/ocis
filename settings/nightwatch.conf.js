const path = require('path')
const WEB_PATH = process.env.WEB_PATH
const TEST_INFRA_DIRECTORY = process.env.TEST_INFRA_DIRECTORY
const OCIS_SETTINGS_STORE = process.env.OCIS_SETTINGS_STORE || '/var/tmp/ocis/settings'

const config = require(path.join(WEB_PATH, 'nightwatch.conf.js'))

config.page_objects_path = [TEST_INFRA_DIRECTORY + '/acceptance/pageObjects', 'ui/tests/acceptance/pageobjects']
config.custom_commands_path = TEST_INFRA_DIRECTORY + '/acceptance/customCommands'

config.test_settings.default.globals = { ...config.test_settings.default.globals, settings_store: OCIS_SETTINGS_STORE }

module.exports = {
  ...config
}
