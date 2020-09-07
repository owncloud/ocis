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
    },

    selectRole: function (username, role) {
      const roleSelector =
        util.format(this.elements.rowByUsername.selector, username) +
        util.format(this.elements.roleInRolesDropdown.selector, role)

      return this
        .click('@rolesDropdownTrigger')
        .waitForElementVisible(roleSelector)
        .click(roleSelector)
    },

    checkUsersRole: function (username, role) {
      const roleSelector =
        util.format(this.elements.rowByUsername.selector, username) +
        util.format(this.elements.currentRole.selector, role)

      return this.useXpath().expect.element(roleSelector).to.be.visible
    },

    toggleUserStatus: function (usernames, status) {
      usernames = usernames.split(',')
      const actionSelector = status === 'enabled' ? this.elements.enableAction : this.elements.disableAction

      // Select users
      for (const username of usernames) {
        const checkboxSelector =
          util.format(this.elements.rowByUsername.selector, username) +
          this.elements.rowCheckbox.selector

        this.useXpath().click(checkboxSelector)
      }

      return this
        .waitForElementVisible('@actionsDropdownTrigger')
        .click('@actionsDropdownTrigger')
        .useCss()
        .waitForElementVisible(actionSelector)
        .click(actionSelector)
        .useXpath()
    },

    checkUsersStatus: function (usernames, status) {
      usernames = usernames.split(',')

      for (const username of usernames) {
        const indicatorSelector =
          util.format(this.elements.rowByUsername.selector, username) +
          util.format(this.elements.statusIndicator.selector, status)

        this.useXpath().waitForElementVisible(indicatorSelector)
      }

      return this
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
    },
    rowByUsername: {
      selector: '//table//td[text()="%s"]/ancestor::tr',
      locateStrategy: 'xpath'
    },
    currentRole: {
      selector: '//span[contains(@class, "accounts-roles-current-role") and normalize-space()="%s"]',
      locateStrategy: 'xpath'
    },
    roleInRolesDropdown: {
      selector: '//label[contains(@class, "accounts-roles-dropdown-role") and normalize-space()="%s"]',
      locateStrategy: 'xpath'
    },
    rolesDropdownTrigger: {
      selector: '//button[contains(@class, "accounts-roles-select-trigger")]',
      locateStrategy: 'xpath'
    },
    loadingAccountsList: {
      selector: '//div[contains(@class, "oc-loader")]',
      locateStrategy: 'xpath'
    },
    rowCheckbox: {
      selector: '//input[@class="oc-checkbox"]',
      locateStrategy: 'xpath'
    },
    actionsDropdownTrigger: {
      selector: '//div[contains(@class, "accounts-actions-dropdown")]//button[normalize-space()="Actions"]',
      locateStrategy: 'xpath'
    },
    disableAction: {
      selector: '#accounts-actions-dropdown-action-disable'
    },
    enableAction: {
      selector: '#accounts-actions-dropdown-action-enable'
    },
    statusIndicator: {
      selector: '//span[contains(@class, "accounts-status-indicator-%s")]',
      locateStrategy: 'xpath'
    }
  }
}
