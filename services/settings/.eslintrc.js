module.exports = {
  "env": {
    "browser": true,
    "es6": true,
    "amd": true
  },
  "extends": [
    "standard",
    "plugin:vue/essential"
  ],
  "parserOptions": {
    "sourceType": "module"
  },
  "rules": {
    'unused-imports/no-unused-imports': 'error'
  },
  plugins: ['unused-imports']
}
