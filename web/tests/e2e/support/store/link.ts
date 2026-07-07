import { Link } from '../types'

export const createdLinkStore = new Map<string, Link>()

export const roleDisplayText: Record<string, string> = {
  'Can view': 'Anyone with the link can view',
  'Can upload': 'Anyone with the link can upload',
  'Can edit': 'Anyone with the link can edit',
  'Secret File Drop': 'Secret File drop'
}
