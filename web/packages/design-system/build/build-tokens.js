import StyleDictionary from 'style-dictionary'
import path, { dirname } from 'node:path'
import yaml from 'yaml'
import jsonFormat from './build-tokens/format-writer-json.js'
import scssFormat from './build-tokens/format-writer-scss.js'
import namespaceTransform from './build-tokens/transform-namespace.js'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const sd = new StyleDictionary({
  hooks: {
    parsers: {
      'yaml-parser': {
        pattern: /\.yaml$/,
        parser: ({ contents }) => yaml.parse(contents)
      }
    }
  },
  parsers: ['yaml-parser'],
  source: [path.join(__dirname, '../src/tokens/**/*.yaml')],
  platforms: {
    ods: {
      transforms: ['name/kebab', 'transform/ods/namespace'],
      buildPath: 'src/assets/tokens/',
      files: [
        {
          destination: 'ods.scss',
          format: 'format/ods/scss',
          filter: ({ filePath }) => filePath.includes('/ods/')
        },
        {
          destination: 'ods.json',
          format: 'format/ods/json',
          filter: ({ filePath }) => filePath.includes('/ods/')
        }
      ]
    }
  },
  log: {
    verbosity: 'verbose'
  }
})

await sd.hasInitialized
sd.registerFormat(jsonFormat)
sd.registerFormat(scssFormat)
sd.registerTransform(namespaceTransform)
await sd.cleanAllPlatforms()
await sd.buildAllPlatforms()
