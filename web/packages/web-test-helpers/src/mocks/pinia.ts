import { createTestingPinia } from '@pinia/testing'
import { mock } from 'vitest-mock-extended'
import {
  AncestorMetaData,
  AppConfigObject,
  ApplicationFileExtension,
  ApplicationInformation,
  ClipboardActions,
  Message,
  Modal,
  OptionsConfig,
  WebThemeType
} from '@ownclouders/web-pkg'
import {
  CollaboratorShare,
  LinkShare,
  Resource,
  ShareRole,
  SpaceResource
} from '@ownclouders/web-client'
import { Group, User } from '@ownclouders/web-client/graph/generated'
import { Capabilities } from '@ownclouders/web-client/ocs'
import defaultTheme from './theme.json'

export { createTestingPinia }

export type PiniaMockOptions = {
  stubActions?: boolean
  appsState?: {
    apps?: Record<string, ApplicationInformation>
    externalAppConfig?: AppConfigObject
    fileExtensions?: ApplicationFileExtension[]
  }
  authState?: {
    accessToken?: string
    idpContextReady?: boolean
    userContextReady?: boolean
    publicLinkContextReady?: boolean
  }
  themeState?: { themes?: WebThemeType[]; currentTheme?: WebThemeType }
  clipboardState?: { action?: ClipboardActions; resources?: Resource[] }
  configState?: {
    server?: string
    options?: OptionsConfig
  }
  messagesState?: { messages?: Message[] }
  modalsState?: { modals?: Modal[] }
  spaceSettingsStore?: {
    spaces?: SpaceResource[]
    selectedSpaces?: SpaceResource[]
  }
  groupSettingsStore?: {
    groups?: Group[]
    selectedGroups?: Group[]
  }
  userSettingsStore?: {
    users?: User[]
    selectedUsers?: User[]
  }
  resourcesStore?: {
    resources?: Resource[]
    currentFolder?: Resource
    ancestorMetaData?: AncestorMetaData
    selectedIds?: string[]
    areFileExtensionsShown?: boolean
    areHiddenFilesShown?: boolean
    deleteQueue?: string[]
  }
  sharesState?: {
    collaboratorShares?: CollaboratorShare[]
    linkShares?: LinkShare[]
    graphRoles?: Record<string, ShareRole>
    loading?: boolean
  }
  spacesState?: { spaces?: SpaceResource[]; currentSpace?: SpaceResource }
  userState?: { user?: User }
  capabilityState?: {
    capabilities?: Partial<Capabilities['capabilities']>
    isInitialized?: boolean
  }
}

export function createMockStore({
  stubActions = true,
  appsState = {},
  authState = {},
  clipboardState = {},
  configState = {},
  themeState = {},
  messagesState = {},
  modalsState = {},
  resourcesStore = {},
  userSettingsStore = {},
  groupSettingsStore = {},
  spaceSettingsStore = {},
  sharesState = {},
  spacesState = {},
  userState = {},
  capabilityState = {}
}: PiniaMockOptions = {}) {
  const defaultOwnCloudTheme = {
    defaults: {
      ...defaultTheme.clients.web.defaults,
      common: {
        ...defaultTheme.common,
        urls: ['https://imprint.url.theme', 'https://privacy.url.theme']
      }
    },
    themes: defaultTheme.clients.web.themes
  }

  return createTestingPinia({
    stubActions,
    initialState: {
      apps: { fileExtensions: [], apps: {}, ...appsState },
      auth: { ...authState },
      clipboard: { resources: [], ...clipboardState },
      config: {
        apps: [],
        external_apps: [],
        customTranslations: [],
        oAuth2: {},
        openIdConnect: {},
        options: {},
        server: '',
        ...configState
      },
      messages: { messages: [], ...messagesState },
      modals: {
        modals: [],
        ...modalsState
      },
      theme: {
        currentTheme: {
          ...defaultOwnCloudTheme.defaults,
          ...defaultOwnCloudTheme.themes[0]
        },
        themes: defaultOwnCloudTheme.themes,
        ...themeState
      },
      resources: { resources: [], ...resourcesStore },
      shares: { collaboratorShares: [], linkShares: [], ...sharesState },
      spaces: { spaces: [], ...spacesState },
      userSettings: { users: [], selectedUsers: [], ...userSettingsStore },
      groupSettings: { groups: [], selectedGroups: [], ...groupSettingsStore },
      spaceSettings: { spaces: [], selectedSpaces: [], ...spaceSettingsStore },
      user: { user: { ...mock<User>({ id: '1' }), ...(userState?.user && { ...userState.user }) } },
      capabilities: {
        graphTagsMaxTagLength: 30,
        ...capabilityState,
        isInitialized: capabilityState?.isInitialized ? capabilityState.isInitialized : true,
        capabilities: {
          ...mock<Capabilities['capabilities']>(),
          ...(capabilityState?.capabilities && { ...capabilityState.capabilities })
        }
      }
    }
  })
}
