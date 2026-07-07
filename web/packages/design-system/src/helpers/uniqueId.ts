let num = 0

export const uniqueId = (prefix: string = '') => {
  prefix = prefix || ''
  num += 1
  return prefix + num
}
