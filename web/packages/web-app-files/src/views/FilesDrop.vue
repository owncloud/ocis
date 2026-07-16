<template>
  <div id="files-drop">
    <app-loading-spinner v-if="loading" />
    <div v-else id="files-drop-wrapper" class="oc-height-1-1">
      <h1 v-if="pageTitle" class="oc-invisible-sr">{{ pageTitle }}</h1>

      <div v-if="dragareaEnabled" class="dragarea" />
      <div key="loaded-drop" class="oc-flex oc-flex-center oc-height-1-1">
        <div class="files-drop-container oc-text-center">
          <div class="oc-width-1-1">
            <h2 v-text="title" />

            <p class="oc-visible@s">
              {{ $gettext("Drop files here to upload or use the 'Upload' button.") }}
            </p>
            <div id="foo">
              <resource-upload
                id="files-drop-zone"
                class="oc-flex oc-flex-middle oc-flex-center oc-placeholder"
              >
                <template #default="{ triggerUpload, uploadLabelId }">
                  <oc-button appearance="filled" variation="primary" @click="triggerUpload">
                    <span
                      :id="uploadLabelId"
                      v-text="
                        $pgettext(
                          'The label of the upload button in secret file upload view',
                          'Upload'
                        )
                      "
                    ></span>
                  </oc-button>
                </template>
              </resource-upload>
            </div>
          </div>

          <div>
            <upload-info
              id="files-drop-upload-info"
              class="oc-width-1-1 oc-height-1-1"
              :info-expanded-initial="true"
              :headless="true"
              :show-expand-details-button="false"
            />
          </div>

          <div>
            <div v-if="errorMessage" class="oc-background-warning oc-width-1-1">
              <h2>
                <span v-text="$gettext('An error occurred while loading the public link')" />
              </h2>
              <p class="oc-rm-m oc-m-rm" v-text="errorMessage" />
            </div>
            <div v-else class="msg-info">
              <oc-icon name="information" />
              <div>
                <p v-text="existingContentNote" />
                <p v-text="flatFolderNote" />
              </div>
            </div>

            <div
              class="oc-width-1-1 oc-flex oc-flex-column oc-flex-middle oc-mt-l"
              :class="{ 'oc-visible@s': isCurrentThemeOwncloud }"
            >
              <oc-img :src="themeLogo" :alt="currentTheme.name" class="oc-visible@s oc-width-1-3" />
              <p class="oc-text-brand-contrast oc-visible@s" v-text="themeSlogan" />

              <template v-if="!isCurrentThemeOwncloud">
                <p
                  class="oc-text-brand-contrast"
                  v-text="$gettext('This feature is brought to you by ownCloud')"
                />
                <oc-button
                  type="a"
                  appearance="raw"
                  :variation="'primary'"
                  size="small"
                  :href="'https://owncloud.com/'"
                  target="_blank"
                >
                  {{ $gettext('Learn more about ownCloud') }}
                </oc-button>
              </template>
            </div>
            <div class="versions mt-5 oc-pb-m oc-pl-m oc-text-xsmall oc-text-muted">
              <span v-text="getWebVersion()" />
              <span v-text="backendVersion" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import {
  createLocationPublic,
  createLocationSpaces,
  useAuthStore,
  useMessages,
  useSpacesStore,
  useThemeStore,
  useUserStore,
  useResourcesStore
} from '@ownclouders/web-pkg'
import ResourceUpload from '../components/AppBar/Upload/ResourceUpload.vue'
import { computed, onMounted, onBeforeUnmount, ref, unref, nextTick, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import {
  useClientService,
  useRouter,
  useRoute,
  useGetMatchingSpace,
  useRouteQuery,
  queryItemAsString,
  useUpload,
  getWebVersion,
  getBackendVersion,
  useCapabilityStore
} from '@ownclouders/web-pkg'
import { eventBus } from '@ownclouders/web-pkg'
import { useService, UppyService } from '@ownclouders/web-pkg'
import { AppLoadingSpinner, useAuthService } from '@ownclouders/web-pkg'
import { HandleUpload } from '../HandleUpload'
import { createFileRouteOptions } from '@ownclouders/web-pkg'
import { PublicSpaceResource, SharePermissionBit } from '@ownclouders/web-client'
import UploadInfo from '../../../web-runtime/src/components/UploadInfo.vue'

const uppyService = useService<UppyService>('$uppyService')
const userStore = useUserStore()
const messageStore = useMessages()
const themeStore = useThemeStore()
const spacesStore = useSpacesStore()
const router = useRouter()
const route = useRoute()
const language = useGettext()
const { $pgettext } = language
const authService = useAuthService()
const clientService = useClientService()
const authStore = useAuthStore()
const { getInternalSpace } = useGetMatchingSpace()
useUpload({ uppyService })

const resourcesStore = useResourcesStore()
const capabilityStore = useCapabilityStore()

const backendVersion = computed(() => getBackendVersion({ capabilityStore }))

const { currentTheme } = storeToRefs(themeStore)

const fileIdQueryItem = useRouteQuery('fileId')

let dragOver: string
let dragOut: string
let drop: string

const share = ref<PublicSpaceResource | null>(null)
const dragareaEnabled = ref(false)
const loading = ref(true)
const errorMessage = ref(null)

const isCurrentThemeOwncloud = computed(() => currentTheme.value.common.name === 'ownCloud')
const themeLogo = computed(() => currentTheme.value.logo.topbar)
const themeSlogan = computed(() => currentTheme.value.common.slogan)

const fileId = computed(() => {
  return queryItemAsString(unref(fileIdQueryItem))
})

const pageTitle = computed(() => {
  return $pgettext(
    'Page title of public link secret file upload view',
    route.value.meta.title as string
  )
})

const title = computed(() => {
  if (unref(share) === null) {
    return ''
  }

  return $pgettext(
    'A message explaining who shared a folder with secret file upload role to the receiving user',
    '%{owner} shared this folder with you for uploading',
    { owner: unref(share).publicLinkShareOwnerDisplayName },
    true
  )
})

const existingContentNote = computed(() => {
  if (unref(share) === null) {
    return ''
  }

  return $pgettext(
    'A note explaining that existing content in secure file drop share is not revealed to anyone else than the owner of the share.',
    'Everyone who has read permission to any parent folder can see the content.',
    { owner: unref(share).publicLinkShareOwnerDisplayName },
    true
  )
})

const flatFolderNote = computed(() => {
  return $pgettext(
    'A note explaining that uploading nested folder structures is not possible in secret file upload.',
    'Transfer of nested folder structures is not possible. Instead, all files from the subfolders will be uploaded individually.'
  )
})

if (!uppyService.getPlugin('HandleUpload')) {
  uppyService.addPlugin(HandleUpload, {
    clientService,
    language,
    route,
    userStore,
    spacesStore,
    messageStore,
    resourcesStore,
    uppyService,
    quotaCheckEnabled: false,
    directoryTreeCreateEnabled: false,
    conflictHandlingEnabled: false
  })
}

const hideDropzone = () => {
  dragareaEnabled.value = false
}
const onDragOver = (event: DragEvent) => {
  dragareaEnabled.value = (event.dataTransfer.types || []).some((e) => e === 'Files')
}

const resolveToInternalLocation = (path: string) => {
  const internalSpace = getInternalSpace(unref(fileId).split('!')[0])
  if (internalSpace) {
    const routeOpts = createFileRouteOptions(internalSpace, { fileId: unref(fileId), path })
    return router.push(createLocationSpaces('files-spaces-generic', routeOpts))
  }

  // no internal space found -> share -> resolve via private link as it holds all the necessary logic
  return router.push({ name: 'resolvePrivateLink', params: { fileId: unref(fileId) } })
}

const resolvePublicLink = async () => {
  loading.value = true

  if (authStore.userContextReady && unref(fileId)) {
    try {
      const path = await clientService.webdav.getPathForFileId(unref(fileId))
      await resolveToInternalLocation(path)
      loading.value = false
      return
    } catch {
      // getPathForFileId failed means the user doesn't have internal access to the resource
    }
  }

  const space = spacesStore.spaces.find(
    (s) => s.driveAlias === `public/${authStore.publicLinkToken}`
  )

  clientService.webdav
    .listFiles(space, {}, { depth: 0 })
    .then(({ resource }) => {
      // Redirect to files list if the link doesn't have role "uploader"
      // FIXME: check for type once https://github.com/owncloud/ocis/issues/8740 is resolved
      const sharePermissions = (resource as PublicSpaceResource).publicLinkPermission
      if (sharePermissions !== SharePermissionBit.Create) {
        router.replace(
          createLocationPublic('files-public-link', {
            params: { driveAliasAndItem: `public/${authStore.publicLinkToken}` }
          })
        )
        return
      }
      share.value = resource as PublicSpaceResource
    })
    .catch((error) => {
      // likely missing password, redirect to public link password prompt
      if (error.statusCode === 401) {
        return authService.handleAuthError(unref(router.currentRoute))
      }
      console.error(error)
      errorMessage.value = error
    })
    .finally(() => {
      loading.value = false
    })
}

watch(loading, async (newLoadValue) => {
  if (!newLoadValue) {
    await nextTick()
    uppyService.useDropTarget({ targetSelector: '#files-drop-wrapper' })
  } else {
    uppyService.removeDropTarget()
  }
})

onMounted(() => {
  dragOver = eventBus.subscribe('drag-over', onDragOver)
  dragOut = eventBus.subscribe('drag-out', hideDropzone)
  drop = eventBus.subscribe('drop', hideDropzone)
  resolvePublicLink()
})

onBeforeUnmount(() => {
  eventBus.unsubscribe('drag-over', dragOver)
  eventBus.unsubscribe('drag-out', dragOut)
  eventBus.unsubscribe('drop', drop)
  uppyService.removeDropTarget()
  uppyService.removePlugin(uppyService.getPlugin('HandleUpload'))
})

// FIXME: remove this once the the test is not interacting with vm directly
defineExpose({ loading })
</script>

<style lang="scss">
#files-drop {
  @media only screen and (min-width: $oc-breakpoint-small-default) {
    padding: var(--oc-space-xlarge);
  }

  #files-drop-wrapper {
    background: transparent;
    padding: var(--oc-space-medium);
    border: none;
    position: relative;
    overflow-y: auto;

    @media (min-width: $oc-breakpoint-small-default) {
      border: 3px dashed var(--oc-color-input-border);
      border-radius: 14px;
      padding: var(--oc-space-xlarge);
    }
  }
  .mt-5 {
    margin-top: 5px;
  }

  &-info-message {
    @media only screen and (min-width: 1200px) {
      width: 400px;
    }
  }

  .msg-info {
    align-items: center;
    background-color: var(--oc-color-background-muted);
    display: flex;
    font-size: var(--oc-font-size-small);
    gap: var(--oc-space-medium);
    padding: var(--oc-space-small);
    text-align: start;

    @media (min-width: $oc-breakpoint-small-default) {
      font-size: var(--oc-font-size-default);
    }
  }
}

.files-drop-container {
  display: grid;
  gap: var(--oc-space-large);
  grid-template-rows: max-content 1fr max-content;
  height: 100%;

  @media (min-width: $oc-breakpoint-xsmall-max) {
    width: 70%;
  }

  @media (min-width: $oc-breakpoint-small-default) {
    width: 55%;
  }

  @media (min-width: $oc-breakpoint-medium-default) {
    width: 30%;
  }
}

.dragarea {
  background-color: rgba(60, 130, 225, 0.21);
  border: none !important;
  pointer-events: none;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  position: absolute;
  z-index: 9;
}
</style>
