<template>
  <div>
    <div
      role="radiogroup"
      class="item-inline-filter oc-flex-inline"
      :class="`item-inline-filter-${filterName}`"
    >
      <oc-button
        v-for="(option, index) in filterOptions"
        :id="option.name"
        :key="index"
        role="radio"
        class="item-inline-filter-option"
        :class="{ 'item-inline-filter-option-selected': activeOption === option.name }"
        :aria-checked="activeOption === option.name"
        appearance="raw"
        @click="toggleFilter(option)"
      >
        <span class="oc-text-truncate item-inline-filter-option-label" v-text="option.label" />
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, unref } from 'vue'
import omit from 'lodash-es/omit'
import { useRoute, useRouteQuery, useRouter, queryItemAsString } from '../../composables'
import { InlineFilterOption } from './types'

interface Props {
  filterName: string
  filterOptions: InlineFilterOption[]
}
interface Emits {
  (e: 'toggleFilter', value: InlineFilterOption): void
}
const { filterName, filterOptions } = defineProps<Props>()
const emit = defineEmits<Emits>()
const router = useRouter()
const currentRoute = useRoute()
const activeOption = ref<string>(filterOptions[0].name)

const queryParam = `q_${filterName}`
const currentRouteQuery = useRouteQuery(queryParam)
const setRouteQuery = (optionName: string) => {
  return router.push({
    query: {
      ...omit(unref(currentRoute).query, [queryParam]),
      [queryParam]: optionName
    }
  })
}

const toggleFilter = async (option: InlineFilterOption) => {
  activeOption.value = option.name
  await setRouteQuery(option.name)
  emit('toggleFilter', option)
}

onMounted(() => {
  const queryStr = queryItemAsString(unref(currentRouteQuery))
  if (queryStr && filterOptions.some(({ name }) => name === queryStr)) {
    activeOption.value = queryStr
    emit(
      'toggleFilter',
      filterOptions.find(({ name }) => name === queryStr)
    )
  }
})
</script>
<style lang="scss">
.item-inline-filter {
  border-radius: 99px;
  border: 1px solid var(--oc-color-text-muted);

  button {
    text-decoration: none;
    font-size: var(--oc-font-size-xsmall);
    line-height: 1rem;
    height: 24px;
    padding: var(--oc-space-xsmall) var(--oc-space-small) !important;
  }

  button:first-child {
    border-top-left-radius: 99px !important;
    border-bottom-left-radius: 99px !important;
    border-top-right-radius: 0px !important;
    border-bottom-right-radius: 0px !important;
  }
  button:last-child {
    border-top-left-radius: 0px !important;
    border-bottom-left-radius: 0px !important;
    border-top-right-radius: 99px !important;
    border-bottom-right-radius: 99px !important;
  }

  &-option-selected {
    background-color: var(--oc-color-swatch-primary-default) !important;
    color: var(--oc-color-text-inverse) !important;
  }
}
</style>
