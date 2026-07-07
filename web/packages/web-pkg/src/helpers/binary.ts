export const blobToArrayBuffer: (blob: Blob) => Promise<string | ArrayBuffer> = (blob: Blob) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onloadend = () => resolve(reader.result)
    reader.onerror = (e) => reject(e)
    reader.readAsArrayBuffer(blob)
  })
}

export const canvasToBlob = (canvas: HTMLCanvasElement): Promise<Blob> => {
  return new Promise((resolve) => canvas.toBlob(resolve))
}
