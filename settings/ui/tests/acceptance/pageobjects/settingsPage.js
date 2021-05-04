const { client } = require('nightwatch-api')
const util = require('util')

module.exports = {
  url: function () {
    return this.api.launchUrl + '/#/settings'
  },

  commands: {
    navigateAndWaitTillLoaded: async function () {
      const url = this.url()
      await await this.navigate(url)
      while (true) {
        let found = false
        await this.waitForElementVisible('@pageHeader', 2000, 500, false)
        await this.api
          .elements('@pageHeader', result => {
            if (result.value.length) {
              found = true
            }
          })
        if (found) {
          break
        }
        await client.refresh()
      }
      return this.waitForElementVisible('@pageHeader')
    },

    getSettingsValue: async function (key) {
      let output
      switch (key) {
        case 'Language':
          let elemfound = true

          // Language value is set to empty at beginning
          // In that case jsut return false
          await this.api.element('@languageValue', result => {
            if (result.status < 0) {
              elemfound = false
            }
          })
          if (!elemfound) {
            output = false
            break
          }
          await this.waitForElementVisible('@languageValue')
            .getText('@languageValue', (result) => {
              output = result.value
            })
          break
        default:
          throw new Error('failed to find the setting')
      }
      return output
    },
    changeSettings: async function (key, value) {
      switch (key) {
        case 'Language':
          await this
            .waitForElementVisible('@languageInput')
            .setValue('@languageInput', value + '\n')
          break
        default:
          throw new Error('failed to find the setting')
      }
    }
  },

  elements: {
    pageHeader: {
      selector: '.oc-page-title'
    },
    languageValue: {
      selector: "//label[.='Language']/..//span[@class='vs__selected']",
      locateStrategy: 'xpath'
    },
    languageInput: {
      selector: "//label[.='Language']/..//input",
      locateStrategy: 'xpath'
    },
  }
}
