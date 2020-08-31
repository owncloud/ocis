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
      const selectXpath = util.format(this.elements.languageSelect.selector, value)
      switch (key) {
        case 'Language':
          await this.waitForElementVisible('@languageValue')
            .click('@languageValue')
            .useXpath()
            .waitForElementVisible(this.elements.languageDropdown.selector)
            .click(selectXpath)
            .waitForElementNotVisible(this.elements.languageDropdown.selector)
            .useCss()
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
      selector: "//label[.='Language']/..//button[starts-with(@id, 'single-choice-toggle')]",
      locateStrategy: 'xpath'
    },
    languageDropdown: {
      selector: "//label[.='Language']/..//div[starts-with(@id, 'single-choice-drop')]",
      locateStrategy: 'xpath'
    },
    languageSelect: {
      selector: "//label[.='Language']/..//div[starts-with(@id, 'single-choice-drop')]//label[normalize-space()='%s']",
      locateStrategy: 'xpath'
    }
  }
}
