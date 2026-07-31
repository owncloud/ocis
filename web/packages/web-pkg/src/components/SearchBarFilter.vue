<template>
  <div class="oc-location-search oc-position-small oc-position-center-right oc-mt-rm" @click.stop>
    <div v-if="currentSelection">
      <oc-filter-chip
        :is-toggle="false"
        :is-toggle-active="false"
        :filter-label="currentSelectionTitle"
        :selected-item-names="[]"
        class="oc-search-bar-filter"
        raw
        close-on-click
      >
        <template #default>
          <oc-button
            v-for="(option, index) in locationOptions"
            :key="index"
            appearance="raw"
            size="medium"
            justify-content="space-between"
            class="search-bar-filter-item oc-flex oc-flex-middle oc-width-1-1 oc-py-xs oc-px-s"
            :class="{ 'oc-mt-s': isIndexGreaterZero(index) }"
            :disabled="!option.enabled"
            :data-test-id="option.id"
            @click="onOptionSelected(option)"
          >
            <oc-icon class="oc-hidden@s" :name="option.icon" />
            <span>{{ option.title }}</span>
            <div v-if="option.id === currentSelection.id" class="oc-flex">
              <oc-icon name="check" />
            </div>
          </oc-button>
        </template>
        <template #active>
          <oc-icon class="oc-hidden@s" :name="currentSelection.icon" />
          <span class="oc-text-truncate oc-invisible-sr@s">{{ currentSelectionTitle }}</span>
        </template>
      </oc-filter-chip>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, Ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import { SearchLocationFilterConstants, useRouteQuery } from '../composables'

type LocationOption = {
  id: string
  title: string
  enabled: Ref<boolean> | boolean
  icon: string
}

interface Props {
  currentFolderAvailable?: boolean
}

interface Emits {
  (e: 'update:modelValue', payload: { value: LocationOption }): void
}

const props = withDefaults(defineProps<Props>(), {
  currentFolderAvailable: false
})
const emit = defineEmits<Emits>()

const { $gettext } = useGettext()
const useScopeQueryValue = useRouteQuery('useScope')

const currentSelection = ref<LocationOption>()
const userSelection = ref<LocationOption>()
const currentSelectionTitle = computed(() => $gettext(currentSelection.value?.title))
const locationOptions = computed<LocationOption[]>(() => [
  {
    id: SearchLocationFilterConstants.currentFolder,
    title: $gettext('Current folder'),
    icon: 'folder',
    enabled: props.currentFolderAvailable
  },
  {
    id: SearchLocationFilterConstants.allFiles,
    title: $gettext('All files'),
    icon: 'globe',
    enabled: true
  }
])

const isIndexGreaterZero = (index: number): boolean => {
  return index > 0
}

watch(
  () => props.currentFolderAvailable,
  () => {
    if (unref(useScopeQueryValue)) {
      const useScope = unref(useScopeQueryValue).toString() === 'true'
      if (useScope) {
        currentSelection.value = unref(locationOptions).find(
          ({ id }) => id === SearchLocationFilterConstants.currentFolder
        )
        return
      }
      currentSelection.value = unref(locationOptions).find(
        ({ id }) => id === SearchLocationFilterConstants.allFiles
      )
      return
    }

    if (!props.currentFolderAvailable) {
      currentSelection.value = unref(locationOptions).find(
        ({ id }) => id === SearchLocationFilterConstants.allFiles
      )
      return
    }

    if (unref(userSelection)) {
      currentSelection.value = unref(locationOptions).find(
        ({ id }) => id === unref(userSelection).id
      )
      return
    }

    currentSelection.value = unref(locationOptions).find(
      ({ id }) => id === SearchLocationFilterConstants.allFiles
    )
  },
  { immediate: true }
)

const onOptionSelected = (option: LocationOption) => {
  userSelection.value = option
  currentSelection.value = option
  emit('update:modelValue', { value: option })
}
</script>

<style lang="scss">
.oc-location-search {
  z-index: 9999;
  margin-right: 34px !important;
  float: right;

  .oc-drop {
    width: 220px;

    @media (min-width: 640px) {
      width: 180px;
    }
  }

  .oc-filter-chip-button {
    justify-content: flex-start;
  }
}
.search-bar-filter-item {
  justify-content: flex-start;

  &:hover {
    background-color: var(--oc-color-background-hover) !important;
  }
}
</style>
