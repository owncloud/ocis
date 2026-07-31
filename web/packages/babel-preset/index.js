const { declare } = require('@babel/helper-plugin-utils')

module.exports = declare(() => {
  return {
    presets: [
      [
        require('@babel/preset-env'),
        {
          useBuiltIns: 'usage',
          shippedProposals: true,
          corejs: {
            version: 3,
            proposals: true
          }
        }
      ]
    ]
  }
})
