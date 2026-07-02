import { Resource } from '@ownclouders/web-client'

export const isSameResource = (r1: Resource, r2: Resource): boolean => {
  if (!r1 || !r2) {
    return false
  }
  return r1.id === r2.id
}
