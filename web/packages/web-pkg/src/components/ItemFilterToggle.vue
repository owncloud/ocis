<template>
  <div class="item-filter oc-flex" :class="`item-filter-${filterName}`">
    <oc-filter-chip
      :is-toggle="true"
      :filter-label="props.filterLabel"
      :is-toggle-active="filterActive"
      @toggle-filter="toggleFilter"
      @clear-filter="toggleFilter"
    />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, unref } from 'vue'
import omit from 'lodash-es/omit'
import { useRoute, useRouteQuery, useRouter, queryItemAsString } from '../composables'

interface Props {
  filterLabel: string
  filterName: string
}

interface Emits {
  (e: 'toggleFilter', isActive: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const router = useRouter()
const currentRoute = useRoute()
const filterActive = ref<boolean>(false)

const queryParam = `q_${props.filterName}`
const currentRouteQuery = useRouteQuery(queryParam)
const setRouteQuery = () => {
  return router.push({
    query: {
      ...omit(unref(currentRoute).query, [queryParam]),
      ...(unref(filterActive) && { [queryParam]: 'true' })
    }
  })
}

const toggleFilter = async () => {
  filterActive.value = !unref(filterActive)
  await setRouteQuery()
  emit('toggleFilter', unref(filterActive))
}

onMounted(() => {
  const queryStr = queryItemAsString(unref(currentRouteQuery))
  if (queryStr === 'true') {
    filterActive.value = true
  }
})
</script>
