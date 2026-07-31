import fs from 'fs'
import path from 'path'
import { config } from '../../config'

// max file creation size is 10GB
export const MAX_FILE_SIZE = Math.pow(1024, 3) * 10

export const getBytes = (fileSize: string): number => {
  fileSize = fileSize.replace(/\s/g, '').toLowerCase()
  const size = parseFloat(fileSize.match(/([\d.]+)/)[0])
  const type = fileSize.match(/[kKmMgGbB]{1,2}/)[0]

  let sizeInbytes = size

  if (!type) {
    return sizeInbytes
  }

  switch (type) {
    case 'b':
      sizeInbytes = size
    case 'kb':
      sizeInbytes = size * 1024
      break
    case 'mb':
      sizeInbytes = size * Math.pow(1024, 2)
      break
    case 'gb':
      sizeInbytes = size * Math.pow(1024, 3)
      break
    default:
      throw new Error('Invalid file size. Must be one of these: b, kb, mb, gb')
  }

  if (sizeInbytes > MAX_FILE_SIZE) {
    throw Error(`File size must be less than '${MAX_FILE_SIZE}' bytes, i.e. 10GB`)
  }

  return sizeInbytes
}

export const getTempUploadPath = (): string => {
  if (!fs.existsSync(config.tempAssetsPath)) {
    fs.mkdirSync(config.tempAssetsPath)
  }
  return config.tempAssetsPath
}

export const createFileWithSize = (
  fileName: string,
  sizeInBytes: number,
  dir: string = getTempUploadPath()
): Promise<void> => {
  return new Promise((resolve, reject) => {
    const fileStream = fs.createWriteStream(path.join(dir, fileName))
    // 500MB buffer size for writing
    const bufferSize = 500 * Math.pow(1024, 2)
    const iterations = Math.ceil(sizeInBytes / bufferSize)

    fileStream.on('open', (fd) => {
      let bytesWritten = 0
      for (let i = 0; i < iterations; i++) {
        const remainingBytes = sizeInBytes - bytesWritten
        const buffer = Buffer.alloc(remainingBytes < bufferSize ? remainingBytes : bufferSize)

        fs.writeSync(fd, buffer as any, 0, buffer.length, null)

        bytesWritten += buffer.length
      }
      fileStream.end()
    })

    fileStream.on('finish', () => resolve())

    fileStream.on('error', (err) => {
      reject(`An error occurred while writing file '${fileName}': ${err}`)
    })
  })
}

export const createFile = (
  fileName: string,
  content: string,
  dir: string = getTempUploadPath()
) => {
  fs.writeFileSync(path.join(dir, fileName), content)
}

export const removeTempUploadDirectory = () => {
  if (fs.existsSync(config.tempAssetsPath)) {
    fs.rmSync(config.tempAssetsPath, { recursive: true })
  }
}
