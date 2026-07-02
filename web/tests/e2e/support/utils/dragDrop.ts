import { readFileSync, lstatSync, readdirSync } from 'fs'
import { Locator, Page } from '@playwright/test'

interface File {
  name: string
  path: string
}

interface FileBuffer {
  name: string
  bufferString: string
  relativePath: string
}

const getFiles = (resources: File[], files: FileBuffer[] = [], parent = '') => {
  for (const resource of resources) {
    const filePath = parent ? `${parent}/${resource.name}` : resource.name
    const stat = lstatSync(resource.path)
    if (stat.isFile()) {
      files.push({
        name: resource.name,
        bufferString: JSON.stringify(Array.from(readFileSync(resource.path))),
        relativePath: filePath
      })
      continue
    }
    // read the directory
    const entries = readdirSync(resource.path)
    const subResources: File[] = entries.map((entry) => ({
      path: resource.path + '/' + entry,
      name: entry
    }))
    getFiles(subResources, files, filePath)
  }
  return files
}

// drag and drop local files to a target element
export const dragDropFiles = async (page: Page, resources: File[], targetSelector: string) => {
  const files = getFiles(resources)

  await page.evaluate<Promise<void>, [FileBuffer[], string]>(
    ([files, targetSelector]) => {
      const dropArea = document.querySelector(targetSelector)
      const dt = new DataTransfer()

      for (const file of files) {
        const buffer = Buffer.from(JSON.parse(file.bufferString))
        const blob = new Blob([buffer])

        const fileObj = new File([blob], file.name)
        // add 'webkitRelativePath' file property only if the file has a parent
        // relative path includes the path with folder structure
        // e.g. folderA/file.txt
        if (file.relativePath.split('/').length > 1) {
          Object.defineProperty(fileObj, 'webkitRelativePath', {
            value: file.relativePath
          })
        }

        dt.items.add(fileObj)
      }

      dropArea.dispatchEvent(new DragEvent('drop', { dataTransfer: dt }))

      return Promise.resolve()
    },
    [files, targetSelector]
  )
}

// drag and drop a element to another element
export const dragTo = async (page: Page, sourceLocator: Locator, destinationLocator: Locator) => {
  // playwright 'dragTo' can be flaky sometimes
  // https://playwright.dev/docs/api/class-locator#locator-drag-to

  // perform drag and drop manually
  await sourceLocator.hover()
  await page.mouse.down()
  await destinationLocator.hover()
  await page.mouse.up()
}
