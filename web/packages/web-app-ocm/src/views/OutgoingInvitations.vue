<template>
  <div class="sciencemesh-app">
    <div>
      <div class="oc-flex oc-flex-middle oc-px-m oc-pt-s">
        <oc-icon name="user-shared" />
        <h2 class="oc-px-s" v-text="$gettext('Invite users to federate')"></h2>
        <oc-contextual-helper class="oc-pl-xs" v-bind="helperContent" />
      </div>
      <div class="oc-flex oc-flex-middle oc-flex-center oc-p-m">
        <oc-button
          :aria-label="
            $gettext('Generate invitation link that can be shared with one or more invitees')
          "
          @click="openInviteModal"
        >
          <oc-icon name="add" />
          <span v-text="$gettext('Generate invitation')" />
        </oc-button>
      </div>
      <oc-modal
        v-if="showInviteModal"
        :title="$gettext('Generate new invitation')"
        :button-cancel-text="$gettext('Cancel')"
        :button-confirm-text="$gettext('Generate')"
        :button-confirm-disabled="!!descriptionErrorMessage"
        focus-trap-initial="#invite_token_description"
        @cancel="resetGenerateInviteToken"
        @confirm="generateToken"
      >
        <template #content>
          <form autocomplete="off" @submit.prevent="generateToken">
            <oc-text-input
              id="invite_token_description"
              v-model="formInput.description"
              class="oc-mb-s"
              :error-message="descriptionErrorMessage"
              :label="$gettext('Add a description (optional)')"
              :clear-button-enabled="true"
              :description-message="
                !descriptionErrorMessage && `${formInput.description?.length || 0}/${50}`
              "
            />
            <input type="submit" class="oc-hidden" />
          </form> </template
      ></oc-modal>
      <app-loading-spinner v-if="loading" />
      <template v-else>
        <no-content-message
          v-if="!sortedTokens.length"
          id="invite-tokens-empty"
          class="files-empty"
          icon="user-shared"
        >
          <template #message>
            <span v-text="$gettext('You have no invitation links')" />
          </template>
        </no-content-message>
        <oc-table
          v-else
          :fields="fields"
          :data="sortedTokens"
          :highlighted="inviteTokensListStore.getLastCreatedToken()"
        >
          <template #token="rowData">
            <div class="invite-code-wrapper oc-flex">
              <div class="oc-text-truncate">
                <span class="oc-text-truncate">{{ encodeInviteToken(rowData.item.token) }}</span>
              </div>
              <oc-button
                id="oc-sciencemesh-copy-token"
                v-oc-tooltip="$gettext('Copy invite token')"
                :aria-label="$gettext('Copy invite token')"
                appearance="raw"
                class="oc-ml-s"
                @click="copyToken(rowData)"
              >
                <oc-icon name="file-copy" />
              </oc-button>
            </div>
          </template>
          <template #link="rowData">
            <a :href="rowData.item.link" v-text="$gettext('Link')" />
            <oc-button
              id="oc-sciencemesh-copy-token"
              v-oc-tooltip="$gettext('Copy invitation link')"
              :aria-label="$gettext('Copy invitation link')"
              appearance="raw"
              @click="copyLink(rowData)"
            >
              <oc-icon name="file-copy" />
            </oc-button>
          </template>
          <template #expiration="rowData">
            <span
              v-oc-tooltip="formatDate(rowData.item.expiration)"
              tabindex="0"
              v-text="formatDateRelative(rowData.item.expiration)"
            />
          </template>
        </oc-table>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, unref } from 'vue'
import {
  NoContentMessage,
  AppLoadingSpinner,
  useClientService,
  useMessages,
  formatDateFromJSDate,
  formatRelativeDateFromJSDate,
  useConfigStore,
  useInviteTokensListStore
} from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { inviteListSchema, inviteSchema } from '../schemas'

const { showMessage, showErrorMessage } = useMessages()
const clientService = useClientService()
const configStore = useConfigStore()
const inviteTokensListStore = useInviteTokensListStore()
const { $gettext, current: currentLanguage } = useGettext()

const showInviteModal = ref(false)
const formInput = ref({
  description: ''
})
const loading = ref(true)
const descriptionErrorMessage = ref<string>()
const fields = computed(() => {
  const haveLinks = unref(sortedTokens)[0]?.link

  return [
    haveLinks && {
      name: 'link',
      title: $gettext('Invitation link'),
      alignH: 'left',
      type: 'slot'
    },
    {
      name: 'token',
      title: $gettext('Invite token'),
      alignH: haveLinks ? 'right' : 'left',
      type: 'slot'
    },
    {
      name: 'description',
      title: $gettext('Description'),
      alignH: 'right'
    },
    {
      name: 'expiration',
      title: $gettext('Expires'),
      alignH: 'right',
      type: 'slot'
    }
  ].filter(Boolean)
})
const sortedTokens = computed(() => {
  return [...unref(inviteTokensListStore.getTokensList())].sort((a, b) =>
    a.expirationSeconds < b.expirationSeconds ? 1 : -1
  )
})
const helperContent = computed(() => {
  return {
    text: $gettext('Create an invitation link and send it to the person you want to share with.'),
    title: $gettext('Invitation link')
  }
})

const encodeInviteToken = (token: string) => {
  const url = new URL(configStore.serverUrl)
  return btoa(`${token}@${url.host}`)
}

const generateToken = async () => {
  const { description } = unref(formInput)

  if (unref(descriptionErrorMessage)) {
    return
  }
  try {
    const { data: tokenInfo } = await clientService.httpAuthenticated.post(
      '/sciencemesh/generate-invite',
      {
        ...(description && { description })
      },
      {
        schema: inviteSchema
      }
    )

    if (tokenInfo.token) {
      inviteTokensListStore.addToken({
        id: tokenInfo.token,
        link: tokenInfo.invite_link,
        token: tokenInfo.token,
        ...(tokenInfo.expiration && {
          expiration: toDateTime(tokenInfo.expiration)
        }),
        ...(tokenInfo.expiration && {
          expirationSeconds: tokenInfo.expiration
        }),
        ...(tokenInfo.description && { description: tokenInfo.description })
      })
      showMessage({
        title: $gettext('Success'),
        status: 'success',
        desc: $gettext(
          'New token has been created and copied to your clipboard. Send it to the invitee(s).'
        )
      })

      const quickToken = encodeInviteToken(tokenInfo.token)
      inviteTokensListStore.setLastCreatedToken(quickToken)
      navigator.clipboard.writeText(quickToken)
    }
  } catch (error) {
    inviteTokensListStore.setLastCreatedToken('')
    errorPopup(error)
  } finally {
    resetGenerateInviteToken()
  }
}

const listTokens = async () => {
  const url = '/sciencemesh/list-invite'
  try {
    const { data } = await clientService.httpAuthenticated.get(url, {
      schema: inviteListSchema
    })
    const tokenList = data.map((t) => ({
      id: t.token,
      token: t.token,
      ...(t.expiration && {
        expiration: toDateTime(t.expiration)
      }),
      ...(t.expiration && {
        expirationSeconds: t.expiration
      }),
      ...(t.description && { description: t.description })
    }))
    inviteTokensListStore.setTokensList(tokenList)
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}

const copyLink = (rowData: { item: { link: string; token: string } }) => {
  navigator.clipboard.writeText(rowData.item.link)
  showMessage({
    title: $gettext('Invition link copied'),
    desc: $gettext('Invitation link has been copied to your clipboard.')
  })
}
const copyToken = (rowData: { item: { link: string; token: string } }) => {
  navigator.clipboard.writeText(encodeInviteToken(rowData.item.token))
  showMessage({
    title: $gettext('Invite token copied'),
    desc: $gettext('Invite token has been copied to your clipboard.')
  })
}
const errorPopup = (error: Error) => {
  console.error(error)
  showErrorMessage({
    title: $gettext('Error'),
    desc: $gettext('An error occurred when generating the token'),
    errors: [error]
  })
}

const openInviteModal = () => {
  showInviteModal.value = true
}

const resetGenerateInviteToken = () => {
  showInviteModal.value = false
  formInput.value = {
    description: ''
  }
}

const toDateTime = (secs: number) => {
  const d = new Date(Date.UTC(1970, 0, 1))
  d.setUTCSeconds(secs)
  return d
}

onMounted(() => {
  listTokens()
})

const formatDate = (date: Date) => {
  return formatDateFromJSDate(date, currentLanguage)
}
const formatDateRelative = (date: Date) => {
  return formatRelativeDateFromJSDate(date, currentLanguage)
}
</script>

<style lang="scss">
.sciencemesh-app {
  .invite-code-wrapper {
    width: 200px;
  }
  #invite-tokens-empty {
    height: 100%;
  }
}
</style>
