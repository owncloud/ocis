import { Page, expect } from '@playwright/test'
import {
  dismissOverlays,
  openSection,
  openSharingPanel,
  selectReportAndOpenSidebar
} from './support/oc'

/**
 * A single captured step. `run` drives the UI into the state to capture; the
 * runner then screenshots it as `output/<tour.id>/NN-<shot>.png`. `title` and
 * `caption` describe the step in the generated `manifest.json`.
 */
export interface Step {
  shot: string
  title: string
  caption: string
  run: (page: Page) => Promise<void>
}

export interface Tour {
  id: string
  source: string
  category: string
  title: string
  summary: string
  steps: Step[]
}

/**
 * Documentation tours derived from the ownCloud Web end-user documentation
 * (doc.owncloud.com/webui/latest/owncloud_web/web_for_users.html).
 *
 * Each tour produces one set of captioned screenshots. Add a tour here, run the
 * capture, and the screenshots plus their captions are produced from the live
 * UI, so the documentation never drifts from the product. New tours need only
 * data plus the `run` actions below.
 */
export const tours: Tour[] = [
  {
    id: 'webtour',
    source: 'doc.owncloud.com · ownCloud Web for users',
    category: 'Getting started',
    title: 'Getting around ownCloud Web',
    summary:
      'A tour of the ownCloud Web interface for everyday use: where your files live, how to switch apps, search everything, and find your account, storage and appearance settings.',
    steps: [
      {
        shot: 'files',
        title: 'Your files',
        caption:
          'After signing in you land in <strong>Personal</strong>, your private space. The left sidebar switches between <strong>Personal</strong>, <strong>Shares</strong>, <strong>Spaces</strong> and <strong>Deleted files</strong>; the main area lists whatever is in the current location.',
        run: async (page) => {
          await expect(page.getByRole('heading', { level: 1, name: 'Personal' })).toBeVisible()
        }
      },
      {
        shot: 'app-switcher',
        title: 'Switch applications',
        caption:
          'The <strong>application switcher</strong> in the top-left corner jumps between Files and the other apps enabled for you, such as the text editor, admin settings or installed extensions.',
        run: async (page) => {
          await page.getByRole('button', { name: 'Application Switcher' }).click()
          await page.waitForTimeout(400)
        }
      },
      {
        shot: 'search',
        title: 'Search everything',
        caption:
          'The <strong>search bar</strong> at the top finds files by name across everything you can access, and by their content when full-text search is enabled. Filters let you narrow a search to a space or file type.',
        run: async (page) => {
          await dismissOverlays(page)
          const search = page.getByRole('combobox', { name: 'Enter search term' })
          await search.click()
          await search.fill('report')
          await page.waitForTimeout(600)
        }
      },
      {
        shot: 'account',
        title: 'Your account and storage',
        caption:
          'The <strong>account menu</strong> in the top-right shows how much storage you have used, opens your <strong>Preferences</strong>, and lets you log out.',
        run: async (page) => {
          const search = page.getByRole('combobox', { name: 'Enter search term' })
          await search.fill('')
          await dismissOverlays(page)
          await page.getByRole('button', { name: 'My Account' }).click()
          await expect(page.getByRole('link', { name: 'Preferences' })).toBeVisible()
        }
      },
      {
        shot: 'preferences',
        title: 'Preferences and appearance',
        caption:
          'Open <strong>Preferences</strong> to set your language, switch between light and dark mode, and review your account details.',
        run: async (page) => {
          await page.goto('/account')
          await expect(page.getByRole('heading', { name: 'Account', level: 1 })).toBeVisible({
            timeout: 30_000
          })
          await page.waitForTimeout(500)
        }
      }
    ]
  },
  {
    id: 'storagetour',
    source: 'doc.owncloud.com · ownCloud Web for users',
    category: 'Getting started',
    title: 'Your files, shares and spaces',
    summary:
      'The left sidebar is your map to everything you can reach: your private files, what others have shared with you, collaborative spaces, and a trash bin for recovering deleted items.',
    steps: [
      {
        shot: 'personal',
        title: 'Personal',
        caption:
          '<strong>Personal</strong> is your private space. Only you can see what is here until you choose to share it. Upload files, create folders and organise your own documents.',
        run: async (page) => {
          await page.goto('/files/spaces/personal')
          await page.waitForURL(/files\/spaces\/personal/, { timeout: 30_000 })
          await expect(page.getByRole('heading', { level: 1, name: 'Personal' })).toBeVisible()
        }
      },
      {
        shot: 'shares',
        title: 'Shared with you',
        caption:
          '<strong>Shares</strong> collects everything other people have shared with you. Depending on the instance you may need to accept a share before it appears, and you can decline ones you do not want.',
        run: async (page) => {
          await openSection(page, 'Shares', /files\/shares/)
        }
      },
      {
        shot: 'spaces',
        title: 'Spaces',
        caption:
          '<strong>Spaces</strong> are shared project areas with their own members and storage. Use them to collaborate with a team in a dedicated location, separate from your personal files.',
        run: async (page) => {
          await openSection(page, 'Spaces', /files\/spaces\/projects/)
        }
      },
      {
        shot: 'deleted',
        title: 'Deleted files',
        caption:
          '<strong>Deleted files</strong> is your trash bin. Restore something you removed by mistake, or permanently delete it to free up space.',
        run: async (page) => {
          await openSection(page, 'Deleted files', /files\/trash/)
        }
      }
    ]
  },
  {
    id: 'filesidebar',
    source: 'doc.owncloud.com · ownCloud Web for users',
    category: 'Getting started',
    title: 'File details, sharing and versions',
    summary:
      'Select any file and the right sidebar shows everything about it: its details, who it is shared with, and its previous versions. This is where most day-to-day file actions happen.',
    steps: [
      {
        shot: 'details',
        title: 'See the details',
        caption:
          'Select a file to open the right sidebar. <strong>Details</strong> shows its size, when it was last modified, the owner, sharing status, tags and how many versions it has.',
        run: async (page) => {
          await selectReportAndOpenSidebar(page)
          await expect(page.getByRole('heading', { level: 3, name: 'report.md' })).toBeVisible({
            timeout: 15_000
          })
        }
      },
      {
        shot: 'share',
        title: 'Share with people',
        caption:
          'The <strong>Shares</strong> panel lets you invite registered users by name, or create a <strong>public link</strong> that anyone can open, optionally protected with a password and an expiry date.',
        run: async (page) => {
          await openSharingPanel(page)
          await page.waitForTimeout(400)
        }
      },
      {
        shot: 'roles',
        title: 'Choose what they can do',
        caption:
          'Pick a role for each person you share with. <strong>Can view</strong> lets them view and download; <strong>Can edit</strong> also lets them upload and change the file. Folders and spaces offer further roles such as Uploader and Manager.',
        run: async (page) => {
          await page.locator('#files-collaborators-role-button-new').click()
          await page.waitForTimeout(500)
        }
      },
      {
        shot: 'versions',
        title: 'Restore an earlier version',
        caption:
          'The <strong>Versions</strong> panel keeps previous copies of a file. Open it to download or restore an older version, which is handy when several people edit the same document.',
        run: async (page) => {
          await page.keyboard.press('Escape')
          await page.waitForTimeout(200)
          // The Shares sub-panel is active; return to the Details panel where the
          // Versions tab lives before opening it.
          const back = page.getByRole('button', { name: 'Back to Details panel' })
          if (await back.isVisible().catch(() => false)) {
            await back.click()
            await page.waitForTimeout(300)
          }
          await page.locator('[data-testid="sidebar-panel-versions-select"]').click()
          await page.waitForTimeout(700)
        }
      }
    ]
  },
  {
    id: 'contextualhelp',
    source: 'doc.owncloud.com · ownCloud Web for users',
    category: 'Getting started',
    title: 'Help where you need it',
    summary:
      'ownCloud Web keeps help close at hand: small question-mark icons sit next to features that benefit from a little explanation, and clicking one shows guidance specific to what you are doing without leaving the page.',
    steps: [
      {
        shot: 'icons',
        title: 'Spot the help icons',
        caption:
          "Throughout the interface a small <strong>?</strong> icon sits next to features that need a little explanation, such as beside <strong>Share with people</strong> and <strong>Public links</strong> in a file's sharing panel.",
        run: async (page) => {
          await selectReportAndOpenSidebar(page)
          await openSharingPanel(page)
          await page.waitForTimeout(300)
        }
      },
      {
        shot: 'popover',
        title: 'Open contextual help',
        caption:
          'Click a help icon to read short, context-specific guidance right where you are, here explaining how sharing with people works, with a <strong>Read more</strong> link to the full documentation.',
        run: async (page) => {
          await page
            .locator('#sidebar-panel-sharing')
            .getByRole('button', { name: 'Show more information' })
            .first()
            .click()
          await page.waitForTimeout(500)
        }
      }
    ]
  }
]
