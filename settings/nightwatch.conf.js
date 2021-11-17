const path = require('path')
const WEB_PATH = process.env.WEB_PATH
const TEST_INFRA_DIRECTORY = process.env.TEST_INFRA_DIRECTORY

const config = require(path.join(TEST_INFRA_DIRECTORY, 'nightwatch.conf.js'))

config.page_objects_path = [TEST_INFRA_DIRECTORY + '/pageObjects', 'ui/tests/acceptance/pageobjects']
config.custom_commands_path = TEST_INFRA_DIRECTORY + '/customCommands'

module.exports = config
