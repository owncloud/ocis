<template>
  <oc-select
    ref="tagSelect"
    v-model="selectedTags"
    class="tags-select"
    :label="$gettext('Tags')"
    :label-hidden="true"
    :multiple="true"
    :disabled="readonly"
    :options="availableTags"
    taggable
    :select-on-key-codes="selectOnKeyCodes"
    :create-option="createOption"
    :selectable="isOptionSelectable"
    :map-keydown="keydownMethods"
    @update:model-value="save"
  >
    <template #selected-option-container="{ option, deselect }">
      <oc-tag class="tags-select-tag oc-ml-xs" :rounded="true" size="small">
        <component
          :is="type"
          v-bind="getAdditionalAttributes(option.label)"
          class="oc-flex oc-flex-middle"
          @click="onTagClicked"
        >
          <oc-icon name="price-tag-3" class="oc-mr-xs" size="small" />
          <span class="oc-text-truncate">{{ option.label }}</span>
        </component>

        <span class="oc-flex oc-flex-middle oc-mr-xs">
          <oc-icon v-if="option.readonly" class="vs__deselect-lock" name="lock" size="small" />
          <oc-button
            v-else
            appearance="raw"
            :title="$gettext('Deselect %{label}', { label: option.label })"
            :aria-label="$gettext('Deselect %{label}', { label: option.label })"
            class="vs__deselect oc-mx-rm"
            @mousedown.stop.prevent
            @click="deselect(option)"
          >
            <oc-icon name="close" size="small" />
          </oc-button>
        </span>
      </oc-tag>
    </template>
    <template #option="{ label, error }">
      <div class="oc-flex test">
        <span class="oc-flex oc-flex-center">
          <oc-tag class="tags-select-tag oc-ml-xs" :rounded="true" size="small">
            <oc-icon name="price-tag-3" size="small" />
            <span class="oc-text-truncate">{{ label }}</span>
          </oc-tag>
        </span>
      </div>
      <div v-if="error" class="oc-text-input-danger">{{ error }}</div>
    </template>
    <template #no-options
      ><span class="oc-text-small oc-text-muted" v-text="$gettext('Enter text to create a Tag')" />
    </template>
  </oc-select>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, toRef, unref, VNodeRef, watch } from 'vue'
import {
  createLocationCommon,
  eventBus,
  SideBarEventTopics,
  useAuthStore,
  useCapabilityStore,
  useClientService,
  useMessages,
  useResourcesStore,
  useRouter
} from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { useTask } from 'vue-concurrency'
import diff from 'lodash-es/difference'
import { call, Resource } from '@ownclouders/web-client'
import { storeToRefs } from 'pinia'

interface Props {
  resource: Resource
}
type TagOption = {
  label: string
  error?: string
  selectable?: boolean
}

const tagsMaxCount = 100

// the keycode property is deprecated in the JS event API, vue-select still works with it though
enum KeyCode {
  Backspace = 8,
  Enter = 13,
  ',' = 188
}

const props = defineProps<Props>()
const { showErrorMessage } = useMessages()
const clientService = useClientService()
const router = useRouter()
const { updateResourceField } = useResourcesStore()

const selectOnKeyCodes = [KeyCode.Enter, KeyCode[',']]

const authStore = useAuthStore()
const { publicLinkContextReady } = storeToRefs(authStore)

const capabilitiesStore = useCapabilityStore()
const { graphTagsMaxTagLength } = storeToRefs(capabilitiesStore)

const type = unref(publicLinkContextReady) ? 'span' : 'router-link'
const resource = toRef(props, 'resource')
const { $gettext } = useGettext()
const readonly = computed(
  () =>
    unref(resource).locked === true ||
    unref(publicLinkContextReady) ||
    (typeof unref(resource).canEditTags === 'function' && unref(resource).canEditTags() === false)
)

const selectedTags = ref<TagOption[]>([])
const availableTags = ref<TagOption[]>([])
let allTags: string[] = []
const tagSelect: VNodeRef = ref(null)

const currentTags = computed<TagOption[]>(() => {
  return [...unref(resource).tags.map((t) => ({ label: t }))]
})

const onTagClicked = () => {
  eventBus.publish(SideBarEventTopics.close)
}

const loadAvailableTagsTask = useTask(function* (signal) {
  const tags = yield* call(clientService.graphAuthenticated.tags.listTags({ signal }))

  allTags = tags
  const selectedLabels = new Set(unref(selectedTags).map((o) => o.label))
  availableTags.value = tags
    .filter((t) => selectedLabels.has(t) === false)
    .map((t) => ({ label: t }))
})

const revertChanges = () => {
  selectedTags.value = unref(currentTags)
}
const createOption = (label: string): TagOption => {
  const len = label.trim().length

  if (!len) {
    return {
      label: label.toLowerCase().trim(),
      error: $gettext('Tag must not consist of blanks only'),
      selectable: false
    }
  }

  if (len > unref(graphTagsMaxTagLength)) {
    return {
      label: label.toLowerCase().trim(),
      error: $gettext('Tags must not be longer than %{max} characters', {
        max: unref(graphTagsMaxTagLength).toString()
      }),
      selectable: false
    }
  }

  return { label: label.toLowerCase().trim() }
}
const isOptionSelectable = (option: TagOption) => {
  return unref(selectedTags).length <= tagsMaxCount && option.selectable !== false
}

const save = async (e: TagOption[] | string[]) => {
  try {
    selectedTags.value = e.map((x) => (typeof x === 'object' ? x : { label: x }))
    const allAvailableTags = new Set([...allTags, ...unref(availableTags).map((t) => t.label)])

    availableTags.value = diff(
      Array.from(allAvailableTags),
      unref(selectedTags).map((o) => o.label)
    ).map((label) => ({
      label
    }))

    const { id, tags, fileId } = unref(resource)
    const selectedTagLabels = unref(selectedTags).map((t) => t.label)
    const tagsToAdd = diff(selectedTagLabels, tags)
    const tagsToRemove = diff(tags, selectedTagLabels)

    if (tagsToAdd.length) {
      await clientService.graphAuthenticated.tags.assignTags({
        resourceId: fileId,
        tags: tagsToAdd
      })
    }

    if (tagsToRemove.length) {
      await clientService.graphAuthenticated.tags.unassignTags({
        resourceId: fileId,
        tags: tagsToRemove
      })
    }

    updateResourceField({ id: id, field: 'tags', value: [...selectedTagLabels] })

    eventBus.publish('sidebar.entity.saved')
    if (unref(tagSelect) !== null) {
      unref(tagSelect).$refs.search.focus()
    }

    allTags.push(...tagsToAdd)
  } catch (e) {
    console.error(e)
    showErrorMessage({
      title: $gettext('Failed to edit tags'),
      errors: [e]
    })
  }
}

watch(resource, () => {
  if (unref(resource)?.tags) {
    revertChanges()
    loadAvailableTagsTask.perform()
  }
})

onMounted(() => {
  if (unref(resource)?.tags) {
    selectedTags.value = unref(currentTags)
  }

  /**
   * If the user can't edit the tags, for example on a public link, there is no need to load the available tags
   */
  if (!unref(readonly)) {
    loadAvailableTagsTask.perform()
  }
})

const keydownMethods = (map: Record<string, (e: Event) => void>) => {
  const objectMapping = {
    ...map
  }
  objectMapping[KeyCode.Backspace] = async (e: Event) => {
    if ((e.target as HTMLInputElement).value || selectedTags.value.length === 0) {
      return
    }

    e.preventDefault()

    availableTags.value.push(selectedTags.value.pop())
    await save(unref(selectedTags))
  }

  return objectMapping
}

const generateTagLink = (tag: string) => {
  const currentTerm = unref(router.currentRoute).query?.term
  return createLocationCommon('files-common-search', {
    query: { provider: 'files.sdk', q_tags: tag, ...(currentTerm && { term: currentTerm }) }
  })
}

const getAdditionalAttributes = (tag: string) => {
  if (unref(publicLinkContextReady)) {
    return {}
  }
  return {
    to: generateTagLink(tag),
    class: 'tags-control-tag-link'
  }
}
</script>

<style lang="scss">
.tags-select {
  .vs__actions {
    display: none !important;
  }

  &-tag {
    height: 1.5rem;

    &-link {
      color: var(--oc-color-swatch-passive-default);
      pointer-events: visible;
    }
  }
}
</style>
