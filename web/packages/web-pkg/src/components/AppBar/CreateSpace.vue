<template>
  <oc-button
    id="new-space-menu-btn"
    key="new-space-menu-btn-enabled"
    v-oc-tooltip="showLabel ? undefined : $gettext('New space')"
    :aria-label="showLabel ? undefined : $gettext('New space')"
    appearance="filled"
    variation="primary"
    @click="showCreateSpaceModal"
  >
    <oc-icon name="add" />
    <span v-if="showLabel" v-text="$gettext('New Space')" />
  </oc-button>
</template>

<script lang="ts" setup>
import { useGettext } from 'vue3-gettext'
import {
  useModals,
  useCreateSpace,
  useSpaceHelpers,
  useMessages,
  useSpacesStore,
  useResourcesStore
} from '../../composables'
import { SpaceResource } from '@ownclouders/web-client'

interface Props {
  showLabel?: boolean
}
interface Emits {
  (event: 'spaceCreated', space: SpaceResource): void
}
const { showLabel = true } = defineProps<Props>()
const emit = defineEmits<Emits>()
const { showMessage, showErrorMessage } = useMessages()
const { $gettext } = useGettext()
const { createSpace } = useCreateSpace()
const { checkSpaceNameModalInput } = useSpaceHelpers()
const { dispatchModal } = useModals()
const spacesStore = useSpacesStore()
const { upsertResource } = useResourcesStore()

const addNewSpace = async (name: string) => {
  try {
    const createdSpace = await createSpace(name, 'project')
    upsertResource(createdSpace)
    spacesStore.upsertSpace(createdSpace)
    emit('spaceCreated', createdSpace)
    showMessage({ title: $gettext('Space was created successfully') })
  } catch (error) {
    console.error(error)
    showErrorMessage({
      title: $gettext('Creating space failed…'),
      errors: [error]
    })
  }
}

const showCreateSpaceModal = () => {
  dispatchModal({
    title: $gettext('Create a new space'),
    confirmText: $gettext('Create'),
    hasInput: true,
    inputLabel: $gettext('Space name'),
    inputValue: $gettext('New space'),
    onConfirm: (name: string) => addNewSpace(name),
    onInput: checkSpaceNameModalInput
  })
}
</script>
