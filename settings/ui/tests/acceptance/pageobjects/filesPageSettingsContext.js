module.exports = {
  commands: {
    getMenuList: async function () {
      const menu = []
      await this.isVisible('@openNavigationBtn', (res) => {
        if (res.value) {
          this.click('@openNavigationBtn')
        }
      })
      await this.waitForElementVisible('@fileSidebarNavItem')
      await this.api
        .elements('@fileSidebarNavItem', result => {
          result.value.map(item => {
            this.api.elementIdText(item.ELEMENT, res => {
              menu.push(res.value)
            })
          })
        })
      return menu
    },
    getUserMenu: async function () {
      const menu = []
      await this
        .waitForElementVisible('@userMenuBtn')
        .click('@userMenuBtn')
        .waitForElementVisible('@userMenuContainer')
      await this.api
        .elements('@userMenuItem', result => {
          result.value.map(item => {
            this.api.elementIdText(item.ELEMENT, res => {
              menu.push(res.value)
            })
          })
        })
      await this
        .waitForElementVisible('@userMenuBtn')
        .click('@userMenuBtn')
        .waitForElementNotVisible('@userMenuContainer')
      return menu
    },
    getFileHeaderItems: async function () {
      const menu = []
      await this.waitForElementVisible('@fileTableHeaderItems')
      await this.api
        .elements('@fileTableHeaderItems', result => {
          result.value.map(item => {
            this.api.elementIdText(item.ELEMENT, res => {
              menu.push(res.value)
            })
          })
        })
      return menu
    }
  },

  elements: {
    pageHeader: {
      selector: '.oc-page-title'
    },
    languageValue: {
      selector: "//button[@id='single-choice-toggle-profile-language']",
      locateStrategy: 'xpath'
    },
    fileSidebarNavItem: {
      selector: '.oc-sidebar-nav-item'
    },
    openNavigationBtn: {
      selector: '//button[@aria-label="Open navigation menu"]',
      locateStrategy: 'xpath'
    },
    userMenuBtn: {
      selector: '#_userMenuButton'
    },
    userMenuItem: {
      selector: '#account-info-container li'
    },
    userMenuContainer: {
      selector: '#account-info-container'
    },
    fileTableHeaderItems: {
      selector: '//*[@id="files-personal-table"]//th[not(.//div)]',
      locateStrategy: 'xpath'
    }
  }
}
