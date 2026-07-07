<template>
  <div class="compare-save-dialog oc-width-1-1 oc-flex oc-flex-between oc-flex-middle">
    <span v-if="saved" class="state-indicator oc-flex oc-flex-middle">
      <oc-icon variation="success" name="checkbox-circle" />
      <span v-translate class="changes-saved oc-ml-s">Changes saved</span>
    </span>
    <span v-else class="state-indicator">{{ unsavedChangesText }}</span>
    <div>
      <oc-button
        :disabled="!unsavedChanges"
        class="compare-save-dialog-revert-btn"
        @click="$emit('revert')"
      >
        <span v-text="$gettext('Revert')" />
      </oc-button>
      <oc-button
        appearance="filled"
        variation="primary"
        class="compare-save-dialog-confirm-btn"
        :disabled="!unsavedChanges || confirmButtonDisabled"
        @click="$emit('confirm')"
      >
        <span v-text="$gettext('Save')" />
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, watch, onBeforeUnmount, onMounted, ref, unref } from 'vue'
import isEqual from 'lodash-es/isEqual'
import { useGettext } from 'vue3-gettext'
import { Group, User } from '@ownclouders/web-client/graph/generated'
import { eventBus } from '../../services/eventBus'

interface Props {
  originalObject: Group | User
  compareObject: Group | User
  confirmButtonDisabled?: boolean
}

interface Emits {
  (e: 'confirm'): void
  (e: 'revert'): void
}

const { originalObject, compareObject, confirmButtonDisabled = false } = defineProps<Props>()
defineEmits<Emits>()

const { $gettext } = useGettext()
const saved = ref(false)
let savedEventToken: string

const unsavedChanges = computed(() => {
  return !isEqual(originalObject, compareObject)
})

const unsavedChangesText = computed(() => {
  return unref(unsavedChanges) ? $gettext('Unsaved changes') : $gettext('No changes')
})

watch(
  () => unsavedChanges.value,
  () => {
    if (unref(unsavedChanges)) {
      saved.value = false
    }
  }
)

watch(
  () => originalObject.id,
  () => {
    saved.value = false
  }
)

onMounted(() => {
  savedEventToken = eventBus.subscribe('sidebar.entity.saved', () => {
    saved.value = true
  })
})

onBeforeUnmount(() => {
  eventBus.unsubscribe('sidebar.entity.saved', savedEventToken)
})
</script>
<style lang="scss" scoped>
.compare-save-dialog {
  background: var(--oc-color-background-highlight);
  flex-flow: row wrap;
}
.state-indicator {
  line-height: 2rem;
}
.changes-saved {
  color: var(--oc-color-swatch-success-default);
}
</style>
