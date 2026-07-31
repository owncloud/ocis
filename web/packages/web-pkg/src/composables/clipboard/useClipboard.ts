import { useClipboard as _useClipboard } from '@vueuse/core'

export const useClipboard = () => {
  // doCopy creates the requested link and copies the url to the clipboard,
  // the copy action uses the clipboard // clipboardItem api to work around the webkit limitations.
  //
  // https://developer.apple.com/forums/thread/691873
  //
  // if those apis not available (or like in firefox behind dom.events.asyncClipboard.clipboardItem)
  // it has a fallback to the vue-use implementation.
  //
  // https://webkit.org/blog/10855/
  const copyToClipboard = (quickLinkUrl: string) => {
    if (typeof ClipboardItem && navigator?.clipboard?.write) {
      return navigator.clipboard.write([
        new ClipboardItem({
          'text/plain': new Blob([quickLinkUrl], { type: 'text/plain' })
        })
      ])
    } else {
      const { copy } = _useClipboard({ legacy: true })
      return copy(quickLinkUrl)
    }
  }

  return { copyToClipboard }
}
