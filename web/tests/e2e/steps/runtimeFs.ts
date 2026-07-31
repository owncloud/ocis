import * as tempFs from '../support/utils/runtimeFs'

export function userCreatesAFileOfSizeInTempUploadDir({
  fileName,
  fileSize
}: {
  fileName: string
  fileSize: string
}) {
  return tempFs.createFileWithSize(fileName, tempFs.getBytes(fileSize))
}
