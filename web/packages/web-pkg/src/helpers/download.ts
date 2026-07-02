export const triggerDownloadWithFilename = (url: string, name: string) => {
  const a = document.createElement('a')
  a.style.display = 'none'
  document.body.appendChild(a)
  a.href = url
  // use download attribute to set desired file name
  a.setAttribute('download', name)
  // trigger the download by simulating click
  a.click()
  // cleanup
  document.body.removeChild(a)
}
