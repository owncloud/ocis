const util = require('util')

module.exports = {
  url: function () {
    return this.api.launchUrl + '/#/accounts'
  },

  commands: {
    navigateAndWaitUntilMounted: async function () {
      const url = this.url()
      return this.navigate(url).waitForElementVisible('@accountsApp')
    },
    accountsList: function () {
      return this.waitForElementVisible('@accountsListTable')
    },
    isUserListed: async function (username) {
      const usernameInTable = util.format(this.elements.userInAccountsList.selector, username)
      await this.useXpath().waitForElementVisible(usernameInTable)
      return true
    },
    isUserDeleted: async function (username) {
      const usernameInTable = util.format(this.elements.userInAccountsList.selector, username)
      await this.useXpath().waitForElementNotPresent(usernameInTable)
      return true
    },

    selectRole: function (username, role) {
      const roleTrigger =
          util.format(this.elements.rowByUsername.selector, username) +
          this.elements.rolesDropdownTrigger.selector
      const roleSelector =
        util.format(this.elements.rowByUsername.selector, username) +
        util.format(this.elements.roleInRolesDropdown.selector, role)

      return this
        .initAjaxCounters()
        .waitForElementVisible(roleTrigger)
        .click(roleTrigger)
        .waitForElementVisible(roleSelector)
        .click(roleSelector)
        .waitForOutstandingAjaxCalls()
    },

    checkUsersRole: function (username, role) {
      const roleSelector =
        util.format(this.elements.rowByUsername.selector, username) +
        util.format(this.elements.currentRole.selector, role)

      return this.useXpath().expect.element(roleSelector).to.be.visible
    },

    setUserActivated: function (usernames, activated) {
      this.selectUsers(usernames)
      return this.click(activated === true ? this.elements.batchActionEnable : this.elements.batchActionDisable)
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
    },

    deleteUsers: function (usernames) {
      this.selectUsers(usernames)
      return this.click(this.elements.batchActionDelete)
        .waitForElementVisible(this.elements.batchActionDeleteConfirm)
        .click(this.elements.batchActionDeleteConfirm)
    },

    selectUsers: function (usernames) {
      usernames = usernames.split(',')

      for (const username of usernames) {
        const checkboxSelector =
          util.format(this.elements.rowByUsername.selector, username) +
          this.elements.rowCheckbox.selector

        this.useXpath().click(checkboxSelector)
      }

      return this
    },

    createUser: function (username, email, password) {
      return this
        .click('@accountsNewAccountTrigger')
        .setValue('@newAccountInputUsername', username)
        .setValue('@newAccountInputEmail', email)
        .setValue('@newAccountInputPassword', password)
        .click('@newAccountButtonConfirm')
    }
  },

  elements: {
    accountsApp: {
      selector: '#accounts-app'
    },
    accountsListTable: {
      selector: '#accounts-user-list'
    },
    userInAccountsList: {
      selector: '//table[@id="accounts-user-list"]//td[text()="%s"]',
      locateStrategy: 'xpath'
    },
    rowByUsername: {
      selector: '//table[@id="accounts-user-list"]//td[text()="%s"]/ancestor::tr',
      locateStrategy: 'xpath'
    },
    currentRole: {
      selector: '//span[contains(@class, "accounts-roles-current-role") and normalize-space()="%s"]',
      locateStrategy: 'xpath'
    },
    roleInRolesDropdown: {
      selector: '//label[contains(@class, "accounts-roles-dropdown-role")]/span[normalize-space()="%s"]',
      locateStrategy: 'xpath'
    },
    rolesDropdownTrigger: {
      selector: '//button[contains(@class, "accounts-roles-select-trigger")]',
      locateStrategy: 'xpath'
    },
    loadingAccountsList: {
      selector: '#accounts-list-loader'
    },
    loadingAccountsListFailed: {
      selector: '#accounts-list-loading-failed'
    },
    rowCheckbox: {
      selector: '//input[contains(@class, "oc-checkbox")]',
      locateStrategy: 'xpath'
    },
    batchActionDisable: {
      selector: '#accounts-batch-action-disable'
    },
    batchActionEnable: {
      selector: '#accounts-batch-action-enable'
    },
    batchActionDelete: {
      selector: '#accounts-batch-action-delete'
    },
    batchActionDeleteCancel: {
      selector: '#accounts-batch-action-delete-cancel'
    },
    batchActionDeleteConfirm: {
      selector: '#accounts-batch-action-delete-confirm'
    },
    statusIndicator: {
      selector: '//span[contains(@class, "accounts-status-indicator-%s")]',
      locateStrategy: 'xpath'
    },
    newAccountInputUsername: {
      selector: '#accounts-new-account-input-username'
    },
    newAccountInputEmail: {
      selector: '#accounts-new-account-input-email'
    },
    newAccountInputPassword: {
      selector: '#accounts-new-account-input-password'
    },
    newAccountButtonConfirm: {
      selector: '#accounts-new-account-button-confirm'
    },
    accountsNewAccountTrigger: {
      selector: '#accounts-new-account-trigger'
    }
  }
}
