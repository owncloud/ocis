import fs from 'fs'
import { User } from '../types'
import { config } from '../../config'

export const userStore = new Map<string, User>([
  [
    'admin',
    {
      id: config.adminUsername,
      displayName: config.adminUsername,
      password: config.adminPassword,
      email: 'admin@example.org'
    }
  ],
  [
    'alice',
    {
      id: 'alice',
      displayName: 'Alice Hansen',
      password: '1234',
      email: 'alice@example.org'
    }
  ],
  [
    'brian',
    {
      id: 'brian',
      displayName: 'Brian Murphy',
      password: '1234',
      email: 'brian@example.org'
    }
  ],
  [
    'carol',
    {
      id: 'carol',
      displayName: 'Carol King',
      password: '1234',
      email: 'carol@example.org'
    }
  ],
  [
    'david',
    {
      id: 'david',
      displayName: 'David William Goodall',
      password: '1234',
      email: 'david@example.org'
    }
  ],
  [
    'edith',
    {
      id: 'edith',
      displayName: 'Edith Anne Widder',
      password: '1234',
      email: 'edith@example.org'
    }
  ],
  [
    'max',
    {
      id: 'max',
      displayName: 'Max Testing',
      password: '12345678',
      email: 'maxtesting@owncloud.com'
    }
  ]
])

export const createdUserStore = new Map<string, User>()

export const federatedUserStore = new Map<string, User>()

// map predefined users to the user store
if (config.predefinedUsers && config.predefinedUsersFile) {
  if (!fs.existsSync(config.predefinedUsersFile)) {
    throw new Error('File not found: ' + config.predefinedUsersFile)
  }
  const users = JSON.parse(fs.readFileSync(config.predefinedUsersFile, 'utf8'))
  for (const [key, user] of Object.entries(users)) {
    userStore.set(key, user as User)
  }
}

// states of the test users:
// - sync enabled/disabled
export const userStateStore = new Map<string, any>()
