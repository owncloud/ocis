import { UsersEnvironment } from '../environment'
import { User } from '../types'

const getValue = (pattern): string => {
  switch (pattern) {
    case '%public%':
      return 'Pwd:12345567'
    default:
      pattern = pattern.replace(/%/g, '')
      const [type, userKey, property] = pattern.split('_')
      if (type === 'user') {
        if (!property) {
          throw new Error('Invalid user property: ' + pattern)
        }
        const usersEnvironment = new UsersEnvironment()
        let user: User
        try {
          user = usersEnvironment.getCreatedUser({ key: userKey })
        } catch (err) {
          // useful for ocm tests where users are from different server
          console.error('[ERR] Failed to get user from created list.', err)
          console.info('[INFO] Getting user from user store.')
          user = usersEnvironment.getUser({ key: userKey })
        }
        return user[property]
      }
  }
}

export const substitute = (text: string): string => {
  if (!text) {
    return text
  }

  const regex = /%[A-Za-z0-9_-]+%/g
  const matches = text.match(regex)
  if (matches) {
    for (const match of matches) {
      const value = getValue(match)
      text = text.replace(match, value)
    }
  }
  return text
}
