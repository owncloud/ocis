import path from 'node:path'
import { color } from 'style-value-types'

export const getPropType = (prop) => {
  const { type } = prop.arguments || {}

  if (type) {
    return type
  } else if (color.test(prop.value)) {
    return 'color'
  } else if (!isNaN(parseInt(prop.value)) || !isNaN(parseFloat(prop.value))) {
    return 'number'
  }

  return '...'
}

export const getPropCategory = (prop) => path.parse(prop.filePath).name

export const sortProps = (props) => {
  return props.sort((a, b) => {
    if (a.name < b.name) {
      return -1
    }
    if (a.name > b.name) {
      return 1
    }

    return 0
  })
}
