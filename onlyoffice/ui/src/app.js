const appInfo = {
  name: 'ocis-onlyoffice',
  id: 'onlyoffice',
  icon: 'x-office-document',
  extensions: [
    {
      extension: 'docx',
      handler: function(config, filePath, fileId) {
        window.open(
          `${config.server}/apps/onlyoffice/${fileId}?filePath=${encodeURIComponent(filePath)}`,
          '_blank'
        )
      },
      newFileMenu: {
        menuTitle($gettext) {
          return $gettext('New OnlyOFFICE document')
        },
        icon: 'x-office-document'
      }
    },
    {
      extension: 'xlsx',
      handler: function(config, filePath, fileId) {
        window.open(
          `${config.server}/apps/onlyoffice/${fileId}?filePath=${encodeURIComponent(filePath)}`,
          '_blank'
        )
      }
    },
    {
      extension: 'pptx',
      handler: function(config, filePath, fileId) {
        window.open(
          `${config.server}/apps/onlyoffice/${fileId}?filePath=${encodeURIComponent(filePath)}`,
          '_blank'
        )
      }
    }
  ]
}

export default {
  appInfo
}

