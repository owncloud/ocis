<template>
  <app-loading-spinner v-if="isLoading" />
  <main v-else id="account" class="oc-mt-m oc-mb-l oc-flex oc-flex-center">
    <div class="account-page">
      <h1 id="account-page-title" v-text="$gettext('My Account')" />
      <account-table
        v-if="showAccountSection"
        :title="$gettext('Account Information')"
        :fields="[$gettext('Information name'), $gettext('Information value')]"
        class="account-page-info"
      >
        <template #header="{ title }">
          <div class="oc-flex oc-flex-between oc-flex-bottom oc-width-1-1">
            <h2 v-text="title" />
            <oc-button
              v-if="accountEditLink"
              variation="primary"
              type="a"
              :href="accountEditLink.href"
              target="_blank"
              data-testid="account-page-edit-url-btn"
            >
              <oc-icon name="edit" />
              <span v-text="$gettext('Edit')" />
            </oc-button>
          </div>
        </template>
        <oc-tr class="account-page-info-username">
          <oc-td>{{ $gettext('Username') }}</oc-td>
          <oc-td translate="no">{{ user.onPremisesSamAccountName }}</oc-td>
        </oc-tr>
        <oc-tr class="account-page-info-displayname">
          <oc-td>{{ $gettext('First and last name') }}</oc-td>
          <oc-td translate="no">{{ user.displayName }}</oc-td>
        </oc-tr>
        <oc-tr class="account-page-info-email">
          <oc-td>{{ $gettext('Email') }}</oc-td>
          <oc-td>
            <template v-if="user.mail">{{ user.mail }}</template>
            <span v-else v-text="$gettext('No email has been set up')" />
          </oc-td>
        </oc-tr>
        <oc-tr
          v-if="user.crossInstanceReference"
          class="account-page-info-cross-instance-reference"
          data-testid="account-cross-instance-reference-row"
        >
          <oc-td>{{
            $pgettext(
              'The oCIS multi-instance reference term displayed in the account page.',
              'Cross-instance reference'
            )
          }}</oc-td>
          <oc-td>
            <span class="cross-instance-reference-wrapper">
              {{ user.crossInstanceReference }}
              <oc-button
                v-oc-tooltip="
                  $pgettext(
                    'The tooltip text for the copy oCIS multi-instance reference button.',
                    'Copy cross-instance reference'
                  )
                "
                appearance="raw"
                variation="passive"
                size="small"
                data-testid="account-cross-instance-reference-copy-btn"
                @click="copyCrossInstanceReference"
              >
                <oc-icon name="file-copy" size="small" />
                <span
                  class="oc-invisible-sr"
                  v-text="
                    $pgettext(
                      'The text of the copy oCIS multi-instance reference button visible only to screen readers.',
                      'Copy cross-instance reference'
                    )
                  "
                />
              </oc-button>
            </span>
          </oc-td>
        </oc-tr>
        <oc-tr v-if="!!quota" class="account-page-info-quota">
          <oc-td>{{ $gettext('Personal storage') }}</oc-td>
          <oc-td data-testid="quota">
            <quota-information :quota="quota" class="oc-mt-xs" />
          </oc-td>
        </oc-tr>
        <oc-tr class="account-page-info-groups">
          <oc-td>{{ $gettext('Group memberships') }}</oc-td>
          <oc-td data-testid="group-names">
            <span v-if="groupNames">{{ groupNames }}</span>
            <span
              v-else
              data-testid="group-names-empty"
              v-text="$gettext('You are not part of any group')"
            />
          </oc-td>
        </oc-tr>
        <oc-tr v-if="showLogout" class="account-page-logout-all-devices">
          <oc-td>{{ $gettext('Logout from active devices') }}</oc-td>
          <oc-td data-testid="logout">
            <oc-button
              appearance="raw"
              type="a"
              :href="logoutUrl"
              target="_blank"
              data-testid="account-page-logout-url-btn"
            >
              <span v-text="$gettext('Show devices')" />
            </oc-button>
          </oc-td>
        </oc-tr>
      </account-table>
      <account-table
        :title="$gettext('Preferences')"
        :fields="[
          $gettext('Preference name'),
          $gettext('Preference description'),
          $gettext('Preference value')
        ]"
        class="account-page-preferences"
      >
        <oc-tr class="account-page-info-language">
          <oc-td>{{ $gettext('Language') }}</oc-td>
          <oc-td>
            <div class="oc-flex">
              <span v-text="$gettext('Select your language.')" />
              <a href="https://explore.transifex.com/owncloud-org/owncloud-web/" target="_blank">
                <div class="oc-flex oc-ml-xs oc-flex-middle">
                  <span v-text="$gettext('Help to translate')" />
                  <oc-icon class="oc-ml-xs" size="small" fill-type="line" name="service" />
                </div>
              </a>
            </div>
          </oc-td>
          <oc-td data-testid="language">
            <oc-select
              v-if="languageOptions"
              :model-value="selectedLanguageValue"
              :label="$gettext('Language')"
              :label-hidden="true"
              :clearable="false"
              :searchable="true"
              :options="languageOptions"
              @update:model-value="updateSelectedLanguage"
            />
          </oc-td>
        </oc-tr>
        <oc-tr v-if="showChangePassword" class="account-page-password">
          <oc-td>{{ $gettext('Password') }}</oc-td>
          <oc-td><span v-text="'**********'" /></oc-td>
          <oc-td data-testid="password">
            <oc-button
              appearance="raw"
              variation="primary"
              data-testid="account-page-edit-password-btn"
              @click="showEditPasswordModal"
            >
              <span v-text="$gettext('Change password')" />
            </oc-button>
          </oc-td>
        </oc-tr>
        <oc-tr class="account-page-info-theme">
          <oc-td>{{ $gettext('Theme') }}</oc-td>
          <oc-td><span v-text="$gettext('Select your favorite theme')" /></oc-td>
          <oc-td data-testid="theme">
            <theme-switcher />
          </oc-td>
        </oc-tr>
        <oc-tr
          v-if="showNotifications && !canConfigureSpecificNotifications"
          class="account-page-notification"
        >
          <oc-td>{{ $gettext('Notifications') }}</oc-td>
          <oc-td v-if="!isMobileWidth">
            <span v-text="$gettext('Receive notification mails')" />
          </oc-td>
          <oc-td data-testid="notification-mails">
            <oc-checkbox
              :model-value="disableEmailNotificationsValue"
              size="large"
              :label="$gettext('Receive notification mails')"
              :label-hidden="!isMobileWidth"
              data-testid="account-page-notification-mails-checkbox"
              @update:model-value="updateDisableEmailNotifications"
            />
          </oc-td>
        </oc-tr>
        <oc-tr v-if="showWebDavDetails" class="account-page-view-options">
          <oc-td>{{ $gettext('View options') }}</oc-td>
          <oc-td v-if="!isMobileWidth">
            <span v-text="$gettext('Show WebDAV information in details view')" />
          </oc-td>
          <oc-td data-testid="view-options">
            <oc-checkbox
              :model-value="viewOptionWebDavDetailsValue"
              size="large"
              :label="$gettext('Show WebDAV information in details view')"
              :label-hidden="!isMobileWidth"
              data-testid="account-page-webdav-details-checkbox"
              @update:model-value="updateViewOptionsWebDavDetails"
            />
          </oc-td>
        </oc-tr>
      </account-table>

      <template v-if="showNotifications && canConfigureSpecificNotifications">
        <account-table
          :title="$gettext('Notifications')"
          :fields="notificationsSettingsFields"
          :show-head="!isMobileWidth"
        >
          <template #header="{ title }">
            <h2>{{ title }}</h2>
            <p>
              {{
                $gettext(
                  'Personalise your notification preferences about any file, folder, or Space.'
                )
              }}
            </p>
          </template>

          <oc-tr v-for="option in notificationsOptions" :key="option.id">
            <oc-td>{{ option.displayName }}</oc-td>
            <oc-td>{{ option.description }}</oc-td>

            <template v-if="option.multiChoiceCollectionValue">
              <oc-td v-for="choice in option.multiChoiceCollectionValue.options" :key="choice.key">
                <span class="checkbox-cell-wrapper">
                  <oc-checkbox
                    :model-value="notificationsValues[option.id][choice.key]"
                    size="large"
                    :label="choice.displayValue"
                    :label-hidden="!isMobileWidth"
                    :disabled="choice.attribute === 'disabled'"
                    @update:model-value="
                      (value) => updateMultiChoiceSettingsValue(option.name, choice.key, value)
                    "
                  />
                </span>
              </oc-td>
            </template>
          </oc-tr>
        </account-table>
        <account-table
          :title="$gettext('Email notification options')"
          :fields="emailNotificationsOptionsFields"
          :show-head="!isMobileWidth"
          class="oc-mt-m"
        >
          <template #header="{ title }">
            <h2 class="oc-invisible-sr">{{ title }}</h2>
          </template>

          <oc-tr v-for="option in emailNotificationsOptions" :key="option.id">
            <oc-td>{{ option.displayName }}</oc-td>
            <oc-td>{{ option.description }}</oc-td>

            <oc-td v-if="option.singleChoiceValue">
              <oc-select
                :model-value="emailNotificationsValues[option.id]"
                :options="option.singleChoiceValue.options"
                :label="emailNotificationsValues[option.id].displayValue"
                :label-hidden="true"
                :clearable="false"
                option-label="displayValue"
                @update:model-value="(value) => updateSingleChoiceValue(option.name, value)"
              />
            </oc-td>
          </oc-tr>
        </account-table>
      </template>

      <account-table
        v-if="extensionPointsWithUserPreferences.length"
        :title="$gettext('Extensions')"
        :fields="[
          $gettext('Extension name'),
          $gettext('Extension description'),
          $gettext('Extension value')
        ]"
        class="account-page-extensions"
      >
        <oc-tr
          v-for="extensionPoint in extensionPointsWithUserPreferences"
          :key="`extension-point-preference-${extensionPoint.id}`"
          class="oc-mb"
        >
          <oc-td>{{ extensionPoint.userPreference.label }}</oc-td>
          <oc-td v-if="extensionPoint.userPreference.description">
            <span v-text="extensionPoint.userPreference.description" />
          </oc-td>
          <oc-td>
            <extension-preference :extension-point="extensionPoint" />
          </oc-td>
        </oc-tr>
      </account-table>
      <account-table
        v-if="showGdprExport"
        :title="$gettext('GDPR')"
        :fields="[
          $gettext('GDPR action name'),
          $gettext('GDPR action description'),
          $gettext('GDPR actions')
        ]"
        class="account-page-gdpr-export"
      >
        <oc-tr class="account-page-gdpr-export">
          <oc-td>{{ $gettext('GDPR export') }}</oc-td>
          <oc-td>
            <span v-text="$gettext('Request a personal data export according to §20 GDPR.')" />
          </oc-td>
          <oc-td data-testid="gdpr-export">
            <gdpr-export />
          </oc-td>
        </oc-tr>
      </account-table>
    </div>
  </main>
</template>

<script lang="ts">
import { storeToRefs } from 'pinia'
import EditPasswordModal from '../components/EditPasswordModal.vue'
import { SettingsBundle, LanguageOption, SettingsValue } from '../helpers/settings'
import { computed, defineComponent, onMounted, onBeforeUnmount, unref, ref } from 'vue'
import {
  useAppsStore,
  useAuthStore,
  useCapabilityStore,
  useClientService,
  useConfigStore,
  useExtensionRegistry,
  useMessages,
  useModals,
  useResourcesStore,
  useSpacesStore,
  useUserStore,
  useSharesStore,
  useClipboard
} from '@ownclouders/web-pkg'
import { useTask } from 'vue-concurrency'
import { useGettext } from 'vue3-gettext'
import { setCurrentLanguage, loadAppTranslations } from '../helpers/language'
import GdprExport from '../components/Account/GdprExport.vue'
import ThemeSwitcher from '../components/Account/ThemeSwitcher.vue'
import ExtensionPreference from '../components/Account/ExtensionPreference.vue'
import { AppLoadingSpinner } from '@ownclouders/web-pkg'
import { SSEAdapter } from '@ownclouders/web-client/sse'
import { supportedLanguages } from '../defaults'
import { User } from '@ownclouders/web-client/graph/generated'
import { isEmpty } from 'lodash-es'
import { call } from '@ownclouders/web-client'
import QuotaInformation from '../components/Account/QuotaInformation.vue'
import AccountTable from '../components/Account/AccountTable.vue'
import { useNotificationsSettings } from '../composables/notificationsSettings'
import { captureException } from '@sentry/vue'

const MOBILE_BREAKPOINT = 800
export default defineComponent({
  name: 'AccountPage',
  components: {
    QuotaInformation,
    AppLoadingSpinner,
    GdprExport,
    ExtensionPreference,
    ThemeSwitcher,
    AccountTable
  },
  setup() {
    const { showMessage, showErrorMessage } = useMessages()
    const userStore = useUserStore()
    const authStore = useAuthStore()
    const language = useGettext()
    const { $gettext, $pgettext } = language
    const clientService = useClientService()
    const resourcesStore = useResourcesStore()
    const appsStore = useAppsStore()
    const sharesStore = useSharesStore()
    const { copyToClipboard } = useClipboard()

    const valuesList = ref<SettingsValue[]>()
    const graphUser = ref<User>()
    const accountBundle = ref<SettingsBundle>()
    const selectedLanguageValue = ref<LanguageOption>()
    const disableEmailNotificationsValue = ref<boolean>()
    const viewOptionWebDavDetailsValue = ref<boolean>(resourcesStore.areWebDavDetailsShown)
    const { dispatchModal } = useModals()
    const spacesStore = useSpacesStore()
    const capabilityStore = useCapabilityStore()
    const configStore = useConfigStore()
    const {
      options: notificationsOptions,
      emailOptions: emailNotificationsOptions,
      values: notificationsValues,
      emailValues: emailNotificationsValues
    } = useNotificationsSettings(valuesList, accountBundle)

    const isMobileWidth = ref<boolean>(window.innerWidth < MOBILE_BREAKPOINT)
    const onResize = () => {
      isMobileWidth.value = window.innerWidth < MOBILE_BREAKPOINT
    }

    // FIXME: Use settings service capability when we have it
    const isSettingsServiceSupported = computed(() => !configStore.options.runningOnEos)

    const { user } = storeToRefs(userStore)

    const quota = computed(() => {
      return spacesStore.personalSpace?.spaceQuota
    })

    const showGdprExport = computed(() => {
      return (
        authStore.userContextReady &&
        capabilityStore.personalDataExport &&
        spacesStore.personalSpace
      )
    })
    const showChangePassword = computed(() => {
      return authStore.userContextReady && !capabilityStore.graphUsersChangeSelfPasswordDisabled
    })
    const showAccountSection = computed(() => authStore.userContextReady)
    const showWebDavDetails = computed(() => authStore.userContextReady)
    const showNotifications = computed(
      () => authStore.userContextReady && unref(isSettingsServiceSupported)
    )
    const showLogout = computed(() => authStore.userContextReady && configStore.options.logoutUrl)

    const loadValuesListTask = useTask(function* (signal) {
      if (!authStore.userContextReady || !unref(isSettingsServiceSupported)) {
        return
      }

      try {
        const {
          data: { values }
        } = yield* call(
          clientService.httpAuthenticated.post<{ values: SettingsValue[] }>(
            '/api/v0/settings/values-list',
            { account_uuid: 'me' },
            { signal }
          )
        )
        valuesList.value = values || []
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Unable to load account data…'),
          errors: [e]
        })
        valuesList.value = []
      }
    }).restartable()

    const loadAccountBundleTask = useTask(function* (signal) {
      if (!authStore.userContextReady || !unref(isSettingsServiceSupported)) {
        return
      }

      try {
        const {
          data: { bundles }
        } = yield* call(
          clientService.httpAuthenticated.post<{ bundles: SettingsBundle[] }>(
            '/api/v0/settings/bundles-list',
            {},
            { signal }
          )
        )
        accountBundle.value = bundles?.find((b) => b.extension === 'ocis-accounts')
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Unable to load account data…'),
          errors: [e]
        })
        accountBundle.value = undefined
      }
    }).restartable()

    const loadGraphUserTask = useTask(function* (signal) {
      if (!authStore.userContextReady) {
        return
      }

      try {
        graphUser.value = yield* call(clientService.graphAuthenticated.users.getMe({}, { signal }))
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Unable to load account data…'),
          errors: [e]
        })
        graphUser.value = undefined
      }
    }).restartable()

    const isLoading = computed(() => {
      return (
        loadValuesListTask.isRunning ||
        !loadValuesListTask.last ||
        loadAccountBundleTask.isRunning ||
        !loadAccountBundleTask.last ||
        loadGraphUserTask.isRunning ||
        !loadGraphUserTask.last
      )
    })

    const languageOptions = Object.keys(supportedLanguages).map((langCode) => ({
      label: supportedLanguages[langCode as keyof typeof supportedLanguages],
      value: langCode
    }))

    const groupNames = computed(() => {
      return unref(user)
        .memberOf.map((group) => group.displayName)
        .join(', ')
    })

    const saveValue = async ({
      identifier,
      valueOptions
    }: {
      identifier: string
      valueOptions: Record<string, any>
    }): Promise<SettingsValue> => {
      let valueId = unref(valuesList)?.find((cV) => cV.identifier.setting === identifier)?.value?.id

      const value = {
        bundleId: unref(accountBundle)?.id,
        settingId: unref(accountBundle)?.settings?.find((s) => s.name === identifier)?.id,
        resource: { type: 'TYPE_USER' },
        accountUuid: 'me',
        ...valueOptions,
        ...(valueId && { id: valueId })
      }

      try {
        const {
          data: { value: data }
        } = await clientService.httpAuthenticated.post<{ value: SettingsValue }>(
          '/api/v0/settings/values-save',
          {
            value: {
              accountUuid: 'me',
              ...value
            }
          }
        )

        // Not sure if we can remove the condition below so just assign this here to be 100% safe
        if (data.value.id) {
          valueId = data.value.id
        }

        /**
         * Edge case: we need to reload the values list to retrieve the valueId if not set,
         * otherwise the backend saves multiple entries
         */
        if (!valueId) {
          loadValuesListTask.perform()
        }

        return data
      } catch (e) {
        throw e
      }
    }

    const updateSelectedLanguage = async (option: LanguageOption) => {
      try {
        loadAppTranslations({
          apps: appsStore.apps,
          gettext: language,
          lang: option.value
        })

        selectedLanguageValue.value = option
        setCurrentLanguage({
          language,
          languageSetting: option.value
        })

        if (authStore.userContextReady) {
          await clientService.graphAuthenticated.users.editMe({
            preferredLanguage: option.value
          } as User)

          /*
           * Refetch shared roles to update their translations when the user changes their preferred language.
           * Otherwise, the roles will remain in the previous language (e.g., English) until the page is refreshed.
           * */
          const shareRolesDefinitions = await sharesStore.fetchShareRolesDefinitions({
            clientService
          })
          sharesStore.setGraphRoles(shareRolesDefinitions)

          if (capabilityStore.supportSSE) {
            ;(clientService.sseAuthenticated as SSEAdapter).updateLanguage(language.current)
          }

          if (spacesStore.personalSpace) {
            // update personal space name with new translation
            spacesStore.updateSpaceField({
              id: spacesStore.personalSpace.id,
              field: 'name',
              value: $gettext('Personal')
            })
          }
        }

        if (loadAccountBundleTask.isRunning) {
          loadAccountBundleTask.cancelAll()
        }

        loadAccountBundleTask.perform()
        showMessage({ title: $gettext('Preference saved.') })
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Unable to save preference…'),
          errors: [e]
        })
      }
    }

    const updateDisableEmailNotifications = async (option: boolean) => {
      try {
        await saveValue({
          identifier: 'disable-email-notifications',
          valueOptions: { boolValue: !option }
        })
        disableEmailNotificationsValue.value = option
        showMessage({ title: $gettext('Preference saved.') })
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Unable to save preference…'),
          errors: [e]
        })
      }
    }

    const updateViewOptionsWebDavDetails = (option: boolean) => {
      try {
        resourcesStore.setAreWebDavDetailsShown(option)
        viewOptionWebDavDetailsValue.value = option
        showMessage({ title: $gettext('Preference saved.') })
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Unable to save preference…'),
          errors: [e]
        })
      }
    }

    const extensionRegistry = useExtensionRegistry()
    const extensionPointsWithUserPreferences = computed(() => {
      return extensionRegistry.getExtensionPoints().filter((extensionPoint) => {
        if (
          !Object.hasOwn(extensionPoint, 'userPreference') ||
          isEmpty(extensionPoint.userPreference)
        ) {
          return false
        }
        const extensions = extensionRegistry.requestExtensions(extensionPoint)
        return !!extensions.length
      })
    })

    const notificationsSettingsFields = computed(() => [
      { label: $gettext('Event') },
      { label: $gettext('Event description'), hidden: true },
      { label: $gettext('In-App'), alignH: 'right' },
      { label: $gettext('Email'), alignH: 'right' }
    ])

    const emailNotificationsOptionsFields = computed(() => [
      { label: $gettext('Options') },
      { label: $gettext('Option description'), hidden: true },
      { label: $gettext('Option value'), hidden: true }
    ])

    const updateValueInValueList = (value: SettingsValue) => {
      const index = unref(valuesList).findIndex(
        (v) => v.identifier.setting === value.identifier.setting
      )

      if (index < 0) {
        valuesList.value.push(value)
        return
      }

      valuesList.value.splice(index, 1, value)
    }

    const updateMultiChoiceSettingsValue = async (
      identifier: string,
      key: string,
      value: boolean | string
    ) => {
      try {
        if (typeof value !== 'boolean') {
          const error = new TypeError(`Unsupported value type ${typeof value}`)

          console.error(error)
          captureException(error)

          return
        }

        const currentValue = unref(valuesList).find((v) => v.identifier.setting === identifier)

        const savedValue = await saveValue({
          identifier,
          valueOptions: {
            collectionValue: {
              values: [
                ...(currentValue?.value.collectionValue.values.filter((val) => val.key !== key) ||
                  []),
                { key, boolValue: value }
              ]
            }
          }
        })

        updateValueInValueList(savedValue)
        showMessage({ title: $gettext('Preference saved.') })
      } catch (error) {
        captureException(error)
        console.error(error)
        showErrorMessage({
          title: $gettext('Unable to save preference…'),
          errors: [error]
        })
      }
    }

    const updateSingleChoiceValue = async (
      identifier: string,
      value: { displayValue: string; value: { stringValue: string } }
    ): Promise<void> => {
      try {
        const savedValue = await saveValue({
          identifier,
          valueOptions: { stringValue: value.value.stringValue }
        })

        updateValueInValueList(savedValue)
        showMessage({ title: $gettext('Preference saved.') })
      } catch (error) {
        captureException(error)
        console.error(error)
        showErrorMessage({
          title: $gettext('Unable to save preference…'),
          errors: [error]
        })
      }
    }

    const canConfigureSpecificNotifications = computed(
      () => capabilityStore.capabilities.notifications.configurable
    )

    onMounted(async () => {
      window.addEventListener('resize', onResize)

      await loadAccountBundleTask.perform()
      await loadValuesListTask.perform()
      await loadGraphUserTask.perform()

      selectedLanguageValue.value = unref(languageOptions)?.find(
        (languageOption) =>
          languageOption.value === (unref(graphUser)?.preferredLanguage || language.current)
      )

      const disableEmailNotificationsConfiguration = unref(valuesList)?.find(
        (cV) => cV.identifier.setting === 'disable-email-notifications'
      )

      disableEmailNotificationsValue.value = disableEmailNotificationsConfiguration
        ? !disableEmailNotificationsConfiguration.value?.boolValue
        : true
    })

    onBeforeUnmount(() => {
      window.removeEventListener('resize', onResize)
    })

    const showEditPasswordModal = () => {
      dispatchModal({
        title: $gettext('Change password'),
        customComponent: EditPasswordModal
      })
    }

    function copyCrossInstanceReference() {
      if (!unref(user).crossInstanceReference) {
        const error = new Error('Attempted to copy cross-instance reference while it is not set')
        console.error(error)
        captureException(error)
        return
      }

      copyToClipboard(unref(user).crossInstanceReference)
      showMessage({
        title: $pgettext(
          'The toast title of the successful copy oCIS multi-instance reference operation.',
          'Cross-instance reference copied'
        ),
        desc: $pgettext(
          'The toast message of the successful copy oCIS multi-instance reference operation.',
          'The cross-instance reference has been copied to your clipboard.'
        )
      })
    }

    return {
      clientService,
      languageOptions,
      extensionPointsWithUserPreferences,
      selectedLanguageValue,
      updateSelectedLanguage,
      updateDisableEmailNotifications,
      updateViewOptionsWebDavDetails,
      accountEditLink: computed(() => configStore.options.accountEditLink),
      showLogout,
      showGdprExport,
      showNotifications,
      showAccountSection,
      showChangePassword,
      showWebDavDetails,
      groupNames,
      user,
      logoutUrl: computed(() => configStore.options.logoutUrl),
      isLoading,
      disableEmailNotificationsValue,
      viewOptionWebDavDetailsValue,
      loadAccountBundleTask,
      loadGraphUserTask,
      loadValuesListTask,
      showEditPasswordModal,
      quota,
      isMobileWidth,
      notificationsOptions,
      notificationsSettingsFields,
      emailNotificationsOptionsFields,
      emailNotificationsOptions,
      notificationsValues,
      updateMultiChoiceSettingsValue,
      emailNotificationsValues,
      updateSingleChoiceValue,
      canConfigureSpecificNotifications,
      copyCrossInstanceReference
    }
  }
})
</script>
<style lang="scss">
#account {
  overflow-y: auto;

  #account-page-title {
    border-bottom: 1px solid var(--oc-color-border);
  }

  .account-page {
    width: 80rem;

    @media (max-width: 1200px) {
      width: 100%;
      padding-left: var(--oc-space-medium);
      padding-right: var(--oc-space-medium);
    }
  }

  .cross-instance-reference-wrapper {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: var(--oc-space-xsmall);
    justify-content: flex-start;
    position: relative;
  }
}
</style>
