import { loadDesignSystem, pages, loadTranslations, supportedLanguages } from './defaults'
import { router } from './router'
import { PortalTarget, useVault } from '@ownclouders/web-pkg'
import { createHead } from '@vueuse/head'
import { abilitiesPlugin } from '@casl/vue'
import { createMongoAbility } from '@casl/ability'

import {
  announceConfiguration,
  initializeApplications,
  announceApplicationsReady,
  announceAuthClient,
  announceDefaults,
  announceClientService,
  announceTheme,
  announcePiniaStores,
  announceCustomStyles,
  announceTranslations,
  announceVersions,
  announceUppyService,
  announceAuthService,
  startSentry,
  announceCustomScripts,
  announceLoadingService,
  announcePreviewService,
  announcePasswordPolicyService,
  registerSSEEventListeners,
  setViewOptions,
  announceGettext,
  announceArchiverService,
  announceAppProviderService
} from './container/bootstrap'
import { applicationStore } from './container/store'
import {
  buildPublicSpaceResource,
  DavHttpError,
  isPersonalSpaceResource,
  isPublicSpaceResource,
  PublicSpaceResource
} from '@ownclouders/web-client'
import { loadCustomTranslations } from './helpers/customTranslations'
import { createApp, onWatcherCleanup, watch } from 'vue'
import PortalVue, { createWormhole } from 'portal-vue'
import { createPinia } from 'pinia'
import Avatar from './components/Avatar.vue'
import focusMixin from './mixins/focusMixin'
import { extensionPoints } from './extensionPoints'
import { isSilentRedirectRoute } from './helpers/silentRedirect'
import { captureException } from '@sentry/vue'
import { CRASH_CODES } from '@ownclouders/web-pkg/src/errors/codes'

export const bootstrapApp = async (configurationPath: string, appsReadyCallback: () => void) => {
  const isSilentRedirect = isSilentRedirectRoute()

  const pinia = createPinia()
  const { isInVault } = useVault()
  const app = createApp(isSilentRedirect ? pages.tokenRenewal : pages.success)
  app.use(pinia)

  const {
    appsStore,
    appStore,
    authStore,
    configStore,
    capabilityStore,
    extensionRegistry,
    spacesStore,
    userStore,
    resourcesStore,
    messagesStore,
    sharesStore,
    webWorkersStore
  } = announcePiniaStores()

  extensionRegistry.registerExtensionPoints(extensionPoints())

  app.provide('$router', router)

  const config = await announceConfiguration({ path: configurationPath, configStore })

  app.use(abilitiesPlugin, createMongoAbility([]), { useGlobalProperties: true })

  const gettext = announceGettext({
    app,
    availableLanguages: supportedLanguages,
    defaultLanguage: config.options.defaultLanguage
  })

  const clientService = announceClientService({ app, configStore, authStore })

  configStore.setIsInVault(isInVault)
  clientService.reinitializeOcsClient(isInVault)
  announceAuthService({
    app,
    configStore,
    router,
    userStore,
    authStore,
    capabilityStore,
    webWorkersStore
  })

  if (!isSilentRedirect) {
    const designSystem = await loadDesignSystem()

    announceUppyService({ app })
    startSentry(configStore, app)
    const appProviderService = announceAppProviderService({
      app,
      serverUrl: configStore.serverUrl,
      clientService
    })
    announceArchiverService({ app, configStore, userStore, capabilityStore })
    announceLoadingService({ app })
    announcePreviewService({
      app,
      configStore,
      userStore,
      authStore,
      capabilityStore
    })
    announcePasswordPolicyService({ app })
    await announceAuthClient(configStore)

    app.config.globalProperties.$wormhole = createWormhole()
    app.use(PortalVue, {
      wormhole: app.config.globalProperties.$wormhole,
      // do not register portal-target component so we can register our own wrapper
      portalTargetName: false
    })
    app.component('PortalTarget', PortalTarget)

    const applicationsPromise = initializeApplications({
      app,
      configStore,
      router,
      appProviderService
    })
    const translationsPromise = loadTranslations()
    const customTranslationsPromise = loadCustomTranslations({ configStore })
    const themePromise = announceTheme({ app, designSystem, configStore })
    const [coreTranslations, customTranslations] = await Promise.all([
      translationsPromise,
      customTranslationsPromise,
      applicationsPromise,
      themePromise
    ])

    // Important: has to happen AFTER native applications are loaded.
    // Reason: the `external` app serves as a blueprint for creating the app provider apps.
    if (applicationStore.has('web-app-external')) {
      await appProviderService.loadData()
      await initializeApplications({
        app,
        configStore,
        router,
        appProviderService,
        appProviderApps: true
      })
    }

    announceTranslations({ appsStore, gettext, coreTranslations, customTranslations })

    announceCustomStyles({ configStore })
    announceCustomScripts({ configStore })
    announceDefaults({ appsStore, router, extensionRegistry, configStore })
  }

  app.use(router)
  app.use(createHead())

  app.component('AvatarImage', Avatar)
  app.mixin(focusMixin)

  app.mount('#owncloud')

  if (isSilentRedirect) {
    return
  }

  setViewOptions({ resourcesStore })

  const applications = Array.from(applicationStore.values())
  applications.forEach((application) => application.mounted(app))

  watch(
    () =>
      authStore.userContextReady || authStore.idpContextReady || authStore.publicLinkContextReady,
    async (newValue, oldValue) => {
      if (!newValue || newValue === oldValue) {
        return
      }
      announceVersions({ capabilityStore })

      await announceApplicationsReady({
        app,
        appsStore,
        applications
      })
      appsReadyCallback()
    },
    {
      immediate: true
    }
  )

  watch(
    () => authStore.userContextReady,
    async (userContextReady) => {
      if (!userContextReady) {
        return
      }

      const clientService = app.config.globalProperties.$clientService
      const previewService = app.config.globalProperties.$previewService
      const passwordPolicyService = app.config.globalProperties.passwordPolicyService
      passwordPolicyService.initialize(capabilityStore)

      // Register SSE event listeners
      if (capabilityStore.supportSSE) {
        registerSSEEventListeners({
          language: gettext,
          resourcesStore,
          spacesStore,
          messageStore: messagesStore,
          sharesStore,
          clientService,
          userStore,
          previewService,
          configStore,
          router
        })
      }

      // load sharing roles from graph API
      const graphRoleDefinitions =
        await clientService.graphAuthenticated.permissions.listRoleDefinitions()
      sharesStore.setGraphRoles(graphRoleDefinitions)

      configStore.setIsInVault(isInVault)
      clientService.reinitializeGraphClient(isInVault)
      clientService.reinitializeOcsClient(isInVault)

      // Load spaces to make them available across the application
      try {
        await spacesStore.loadSpaces({ graphClient: clientService.graphAuthenticated, isInVault })
        const personalSpace = spacesStore.spaces.find(isPersonalSpaceResource)

        if (personalSpace) {
          spacesStore.updateSpaceField({
            id: personalSpace.id,
            field: 'name',
            value: app.config.globalProperties.$gettext('Personal')
          })
        }
      } catch (error) {
        console.error(error)
        captureException(error)
        router.push({ name: 'crash', query: { code: CRASH_CODES.RUNTIME_BOOTSTRAP_SPACES_LOAD } })
      }
    },
    {
      immediate: true
    }
  )
  watch(
    () => authStore.publicLinkContextReady,
    async (publicLinkContextReady) => {
      appStore.error = null
      if (!publicLinkContextReady) {
        return
      }
      // Create virtual space for public link
      const publicLinkToken = authStore.publicLinkToken
      const publicLinkPassword = authStore.publicLinkPassword
      const publicLinkType = authStore.publicLinkType

      const space = buildPublicSpaceResource({
        id: publicLinkToken,
        name: app.config.globalProperties.$gettext('Public files'),
        ...(publicLinkPassword && { publicLinkPassword }),
        serverUrl: configStore.serverUrl,
        publicLinkType: publicLinkType
      })

      spacesStore.addSpaces([space])

      const controller = new AbortController()

      onWatcherCleanup(() => {
        controller.abort()
      })

      try {
        // Return early if the password is required. Let the resolvePublicLink.`verifyPassword` handle that.
        if (authStore.publicLinkPasswordRequired) {
          return
        }
        const loadedSpace = await clientService.webdav.getFileInfo(
          space,
          {},
          { signal: controller.signal }
        )

        for (const key in loadedSpace) {
          if (loadedSpace[key] !== undefined) {
            space[key] = loadedSpace[key]
          }
        }

        spacesStore.upsertSpace(space)
      } catch (error) {
        const err = error as DavHttpError

        if (err.statusCode === 401) {
          return
        }

        if (err.statusCode === 404) {
          const notFoundError = new Error(
            app.config.globalProperties.$gettext(
              'The resource could not be located, it may not exist anymore.'
            )
          )
          appStore.error = notFoundError
          throw notFoundError
        }

        throw err
      } finally {
        spacesStore.setSpacesInitialized(true)
      }
    },
    {
      immediate: true
    }
  )
  watch(
    // only needed if a public link gets re-resolved with a changed password prop (changed or removed).
    // don't need to set { immediate: true } on the watcher.
    () => authStore.publicLinkPassword,
    (publicLinkPassword: string | undefined) => {
      const publicLinkToken = authStore.publicLinkToken
      const space = spacesStore.spaces.find((space) => {
        return isPublicSpaceResource(space) && space.id === publicLinkToken
      })
      if (!space) {
        return
      }
      ;(space as PublicSpaceResource).publicLinkPassword = publicLinkPassword
    }
  )
}

export const bootstrapErrorApp = async (err: Error): Promise<void> => {
  const useBrowserLanguage = (language: string): string => {
    const fallbackLanguage = 'en'

    return supportedLanguages[language] ? language : fallbackLanguage
  }
  const { capabilityStore, configStore } = announcePiniaStores()
  announceVersions({ capabilityStore })
  const app = createApp(pages.failure)
  const designSystem = await loadDesignSystem()
  await announceTheme({ app, designSystem, configStore })
  console.error(err)
  const translations = await loadTranslations()
  const gettext = announceGettext({
    app,
    availableLanguages: supportedLanguages,
    defaultLanguage: useBrowserLanguage(window.navigator.language.slice(0, 2))
  })
  announceTranslations({ gettext, coreTranslations: translations })
  app.mount('#owncloud')
}
;(window as any).runtimeLoaded({
  bootstrapApp,
  bootstrapErrorApp
})
