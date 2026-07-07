import { Group } from '../types'

export const dummyGroupStore = new Map<string, Group>([
  [
    'security',
    {
      id: 'security',
      displayName: 'security department'
    }
  ],
  [
    'sales',
    {
      id: 'sales',
      displayName: 'sales department'
    }
  ],
  [
    'finance',
    {
      id: 'finance',
      displayName: 'finance department'
    }
  ]
])

export const createdGroupStore = new Map<string, Group>()
