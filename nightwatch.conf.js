const path = require('path')
const PHOENIX_PATH = process.env.PHOENIX_PATH

const config = require(path.join(PHOENIX_PATH, 'nightwatch.conf.js'))

config.page_objects_path = [PHOENIX_PATH + '/tests/acceptance/pageObjects', 'ui/tests/acceptance/pageobjects']
config.custom_commands_path = PHOENIX_PATH + '/tests/acceptance/customCommands'

module.exports = {
  ...config
}
