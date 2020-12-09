const path = require('path')
const WEB_PATH = process.env.WEB_PATH

const config = require(path.join(WEB_PATH, 'nightwatch.conf.js'))

config.page_objects_path = [WEB_PATH + '/tests/acceptance/pageObjects', 'ui/tests/acceptance/pageobjects']
config.custom_commands_path = WEB_PATH + '/tests/acceptance/customCommands'

module.exports = {
  ...config
}
