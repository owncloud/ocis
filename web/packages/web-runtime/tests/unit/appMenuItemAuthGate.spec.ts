import { Ref, unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'
import { ApplicationInformation, Extension, useUserStore } from '@ownclouders/web-pkg'
import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'

import activities from '../../../web-app-activities/src/index'
import adminSettings from '../../../web-app-admin-settings/src/index'
import appStore from '../../../web-app-app-store/src/index'
import epubReader from '../../../web-app-epub-reader/src/index'
import textEditor from '../../../web-app-text-editor/src/index'
import { extensions as filesExtensions } from '../../../web-app-files/src/extensions'
import { extensions as ocmExtensions } from '../../../web-app-ocm/src/extensions'

/**
 * Every app that contributes an entry to the top bar app drawer must withhold it
 * from anonymous visitors. The app drawer renders inside the full application
 * layout, which a public link file view also gets, and the extension registry is
 * shared between the authenticated and the public link context - so an
 * unguarded registration is visible to anyone holding a link.
 *
 * The apps expose their extensions in two shapes: most return them from
 * `setup()`, while files and ocm export an `extensions(appInfo)` composable.
 */
const APPS: Array<{ name: string; getExtensions: () => Ref<Extension[]> }> = [
  {
    name: 'activities',
    getExtensions: () => activities.setup({ applicationConfig: {} }).extensions
  },
  {
    name: 'admin-settings',
    getExtensions: () => adminSettings.setup({ applicationConfig: {} }).extensions
  },
  { name: 'app-store', getExtensions: () => appStore.setup({ applicationConfig: {} }).extensions },
  {
    name: 'epub-reader',
    getExtensions: () => epubReader.setup({ applicationConfig: {} }).extensions
  },
  {
    name: 'text-editor',
    getExtensions: () => textEditor.setup({ applicationConfig: {} }).extensions
  },
  {
    name: 'files',
    getExtensions: () => filesExtensions(mock<ApplicationInformation>({ id: 'files' }))
  },
  { name: 'ocm', getExtensions: () => ocmExtensions(mock<ApplicationInformation>({ id: 'ocm' })) }
]

/**
 * Apps gated on nothing but a signed-in user. admin-settings and app-store are
 * absent on purpose: they additionally require permissions, so "a user is signed
 * in" is not enough to expect their menu item.
 */
const USER_ONLY_GATED = ['activities', 'epub-reader', 'files', 'ocm', 'text-editor']

const getAppMenuItems = (
  getExtensions: () => Ref<Extension[]>,
  { user }: { user: User | null }
) => {
  let items: Extension[]
  // Some apps reach for the router while building their extensions, so the
  // wrapper has to supply one.
  const mocks = defaultComponentMocks()
  getComposableWrapper(
    () => {
      // The shared pinia mock always seeds a user, so an anonymous visitor has to
      // be expressed by writing the store directly.
      const userStore = useUserStore()
      userStore.user = user

      items = unref(getExtensions()).filter((e) => e.type === 'appMenuItem')
    },
    { mocks, provide: mocks }
  )
  return items
}

describe('app menu item auth gate', () => {
  it.each(APPS)(
    'the $name app contributes no app menu item for an anonymous visitor',
    ({ getExtensions }) => {
      expect(getAppMenuItems(getExtensions, { user: null })).toHaveLength(0)
    }
  )
  // The counterpart: the gate must not swallow the menu item for real users.
  it.each(APPS.filter(({ name }) => USER_ONLY_GATED.includes(name)))(
    'the $name app contributes its app menu item for a signed-in user',
    ({ getExtensions }) => {
      expect(getAppMenuItems(getExtensions, { user: mock<User>({ id: '1' }) })).toHaveLength(1)
    }
  )
})
