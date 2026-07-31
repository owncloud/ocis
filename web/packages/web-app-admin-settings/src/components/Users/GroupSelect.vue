<template>
  <div id="user-group-select-form">
    <oc-select
      :model-value="selectedOptions"
      class="oc-mb-s"
      :multiple="true"
      :options="groupOptions"
      option-label="displayName"
      :label="$gettext('Groups')"
      :fix-message-line="true"
      v-bind="$attrs"
      @update:model-value="onUpdate"
    >
      <template #selected-option="{ displayName, id }">
        <span class="oc-flex oc-flex-center">
          <avatar-image
            class="oc-flex oc-align-self-center oc-mr-s"
            :width="16.8"
            :userid="id"
            :user-name="displayName"
          />
          <span>{{ displayName }}</span>
        </span>
      </template>
      <template #option="{ displayName, id }">
        <div class="oc-flex">
          <span class="oc-flex oc-flex-center">
            <avatar-image
              class="oc-flex oc-align-self-center oc-mr-s"
              :width="16.8"
              :userid="id"
              :user-name="displayName"
            />
            <span>{{ displayName }}</span>
          </span>
        </div>
      </template>
    </oc-select>
  </div>
</template>
<script lang="ts" setup>
import { computed, ref, unref, watch } from 'vue'
import { Group } from '@ownclouders/web-client/graph/generated'

interface Props {
  selectedGroups: Group[]
  groupOptions: Group[]
}
interface Emits {
  (e: 'selectedOptionChange', value: Group): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const selectedOptions = ref()
const onUpdate = (group: Group) => {
  selectedOptions.value = group
  emit('selectedOptionChange', unref(selectedOptions))
}

const currentGroups = computed(() => props.selectedGroups)
watch(
  currentGroups,
  () => {
    selectedOptions.value = props.selectedGroups
      .map((g) => ({
        ...g,
        readonly: g.groupTypes?.includes('ReadOnly')
      }))
      .sort((a: any, b: any) => b.readonly - a.readonly)
  },
  { immediate: true }
)
</script>
