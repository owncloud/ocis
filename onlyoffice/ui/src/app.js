const appInfo = {
  name: 'OnlyOffice',
  id: 'onlyoffice',
  icon: 'x-office-document',
  extensions: [
    {
      extension: 'docx',
      handler: function({ extensionConfig, filePath, fileId }) {
        window.open(
          `${extensionConfig.server}/apps/onlyoffice/${fileId}?filePath=${encodeURIComponent(filePath)}`,
          '_blank'
        )
      },
      newFileMenu: {
        menuTitle($gettext) {
          return $gettext('New OnlyOffice document')
        },
        icon: 'x-office-document'
      }
    },
    {
      extension: 'xlsx',
      handler: function({ extensionConfig, filePath, fileId }) {
        window.open(
          `${extensionConfig.server}/apps/onlyoffice/${fileId}?filePath=${encodeURIComponent(filePath)}`,
          '_blank'
        )
      }
    },
    {
      extension: 'pptx',
      handler: function({ extensionConfig, filePath, fileId }) {
        window.open(
          `${extensionConfig.server}/apps/onlyoffice/${fileId}?filePath=${encodeURIComponent(filePath)}`,
          '_blank'
        )
      }
    }
  ]
}

export default {
  appInfo
}

