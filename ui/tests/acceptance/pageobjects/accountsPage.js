const util = require('util')

module.exports = {
  url: function () {
    return this.api.launchUrl + '/#/accounts'
  },

  commands: {
    navigateAndWaitTillLoaded: async function () {
      const url = this.url()
      return this.navigate(url).waitForElementVisible('@accountsLabel')
    },
    accountsList: function () {
      return this.waitForElementVisible('@accountsListTable')
    },
    isUserListed: async function (username) {
      let user
      const usernameInTable = util.format(this.elements.userInAccountsList.selector, username)
      await this.useXpath().waitForElementVisible(usernameInTable)
        .getText(usernameInTable, (result) => {
          user = result
        })
      return user.value
    }
  },

  elements: {
    accountsLabel: {
      selector: "//h1[normalize-space(.)='Accounts']",
      locateStrategy: 'xpath'
    },
    accountsListTable: {
      selector: "//table[@class='uk-table uk-table-middle uk-table-divider']",
      locateStrategy: 'xpath'
    },
    userInAccountsList: {
      selector: '//table//td[text()="%s"]',
      locateStrategy: 'xpath'
    }
  }
}
