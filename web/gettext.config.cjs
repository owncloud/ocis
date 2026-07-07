module.exports = {
  input: {
    path: './src',
    include: ['**/*.js', '**/*.ts', '**/*.vue']
  },
  output: {
    locales: [
      'af',
      'ar',
      'bs',
      'bg',
      'ca',
      'cs',
      'de',
      'el',
      'es',
      'et',
      'fr',
      'gl',
      'he',
      'hr',
      'id',
      'it',
      'ja',
      'nl',
      'pl',
      'pt',
      'ka',
      'ko',
      'ro',
      'ru',
      'si',
      'sk',
      'sq',
      'sv',
      'sr',
      'ta',
      'tr',
      'ug',
      'uk',
      'zh'
    ],
    path: './l10n/locale',
    potPath: '../template.pot',
    jsonPath: '../translations.json',
    flat: false,
    linguas: false
  }
}
