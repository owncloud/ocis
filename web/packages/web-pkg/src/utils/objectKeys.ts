import isPlainObject from 'lodash-es/isPlainObject'

export const objectKeys = (obj: Record<string, unknown>): string[] => {
  const paths: string[] = []

  const walk = (o: Record<string, unknown>, p = '') =>
    Object.keys(o).forEach((key) => {
      if (isPlainObject(o[key])) {
        walk(o[key] as Record<string, unknown>, `${p}${key}.`)
      } else {
        paths.push(`${p}${key}`)
      }
    })

  walk(obj)

  return paths
}
