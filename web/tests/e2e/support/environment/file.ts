import fs from 'fs'
import path from 'path'
import { config } from '../../config'
import { File } from '../types'

export class FilesEnvironment {
  getFile({ name }: { name: string }): File {
    const relPath = path.join(config.assets, name)
    if (!fs.existsSync(relPath)) {
      throw new Error('TODO: fixture files')
    }

    return { name, path: path.resolve(relPath) }
  }
}
