import { APIRequestContext } from '@playwright/test'
import { config } from '../../config'

const USER = config.adminUsername
const PASS = config.adminPassword

const basic = (user: string, pass: string) =>
  'Basic ' + Buffer.from(`${user}:${pass}`).toString('base64')

/**
 * Best-effort, idempotent demo data so the tours have something to show:
 * a versioned file, a trashed item, a project space, and one incoming share.
 * Every step is independent and tolerates "already exists" responses, so the
 * tool is safe to run repeatedly. Failures are logged, never fatal — the tours
 * simply capture whatever state the instance is in.
 */
export async function seed(api: APIRequestContext): Promise<void> {
  const admin = basic(USER, PASS)

  // A versioned text file (details / versions / sharing / contextual-help tours).
  try {
    const versions = [
      '# Project report\n\nInitial draft of the quarterly report.\n',
      '# Project report\n\nQuarterly report with the Q3 budget and an updated timeline.\n'
    ]
    for (const data of versions) {
      await api.put(`/remote.php/dav/files/${USER}/report.md`, {
        headers: { Authorization: admin, 'Content-Type': 'text/markdown' },
        data
      })
      await new Promise((resolve) => setTimeout(resolve, 1100))
    }
  } catch (error) {
    console.warn('[seed] report.md:', error)
  }

  // A trashed item (deleted-files tour).
  try {
    await api.put(`/remote.php/dav/files/${USER}/old-draft.txt`, {
      headers: { Authorization: admin },
      data: 'An old draft, later deleted.'
    })
    await api.delete(`/remote.php/dav/files/${USER}/old-draft.txt`, {
      headers: { Authorization: admin }
    })
  } catch (error) {
    console.warn('[seed] trash:', error)
  }

  // A project space (spaces tour).
  try {
    await api.post('/graph/v1.0/drives', {
      headers: { Authorization: admin, 'Content-Type': 'application/json' },
      data: { name: 'Demo project', driveType: 'project' }
    })
  } catch (error) {
    console.warn('[seed] project space:', error)
  }

  // An incoming share (shared-with-you tour): a helper user shares a file with the admin.
  try {
    const helper = 'doc-demo'
    const helperPass = 'Secret123!'
    await api.post('/graph/v1.0/users', {
      headers: { Authorization: admin, 'Content-Type': 'application/json' },
      data: {
        displayName: 'Documentation Demo',
        onPremisesSamAccountName: helper,
        mail: 'doc-demo@example.org',
        accountEnabled: true,
        passwordProfile: { password: helperPass }
      }
    })
    const helperAuth = basic(helper, helperPass)
    await api.put(`/remote.php/dav/files/${helper}/shared-notes.txt`, {
      headers: { Authorization: helperAuth },
      data: 'A file shared with you.'
    })
    await api.post('/ocs/v2.php/apps/files_sharing/api/v1/shares?format=json', {
      headers: { Authorization: helperAuth, 'OCS-APIRequest': 'true' },
      form: { path: '/shared-notes.txt', shareType: '0', shareWith: USER, permissions: '1' }
    })
  } catch (error) {
    console.warn('[seed] incoming share:', error)
  }
}
