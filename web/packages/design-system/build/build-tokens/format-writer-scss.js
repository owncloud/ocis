import { sortProps } from './utils.js'

export default {
  name: 'format/ods/scss',
  format: (dictionary) => {
    const props = sortProps(dictionary.allTokens)
    const data = [
      ...props.map((p) => `$${p.name}: ${p.value};`),
      '',
      ':host, :root {',
      ...props.map((p) => `  --${p.name}: #{$${p.name}};`),
      '}',
      ''
    ].join('\n')

    return data
  }
}
