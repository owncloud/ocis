<template>
  <div class="date-filter oc-flex" :class="`date-filter-${filterName}`">
    <oc-filter-chip
      ref="filterChip"
      :filter-label="filterLabel"
      :selected-item-names="selectedItem ? [selectedItem[displayNameAttribute]] : undefined"
      @clear-filter="clearFilter"
      @show-drop="onShowDrop"
      @hide-drop="onHideDrop"
    >
      <template #default>
        <oc-list class="date-filter-list" :class="{ 'date-filter-list-hidden': dateRangeClicked }">
          <li v-for="(item, index) in displayedItems" :key="index" class="oc-my-xs">
            <oc-button
              class="date-filter-list-item oc-flex oc-flex-between oc-flex-middle oc-width-1-1 oc-p-xs"
              :class="{
                'date-filter-list-item-active': isItemSelected(item)
              }"
              justify-content="space-between"
              appearance="raw"
              :data-testid="item[displayNameAttribute]"
              @click="toggleItemSelection(item)"
            >
              <div class="oc-flex oc-flex-middle oc-text-truncate">
                <div class="oc-text-truncate oc-ml-s">
                  <slot name="item" :item="item" />
                </div>
              </div>
              <div class="oc-flex">
                <oc-icon v-if="isItemSelected(item)" name="check" />
              </div>
            </oc-button>
          </li>
          <li class="oc-my-xs">
            <oc-button
              class="date-filter-list-item oc-flex oc-flex-between oc-flex-middle oc-width-1-1 oc-p-xs"
              :class="{
                'date-filter-list-item-active': dateRangeApplied
              }"
              justify-content="space-between"
              appearance="raw"
              data-testid="custom-date-range"
              @click="dateRangeClicked = true"
            >
              <div class="oc-flex oc-flex-middle oc-text-truncate">
                <div class="oc-text-truncate oc-ml-s">
                  <span v-text="$gettext('Custom date range')" />
                </div>
              </div>
              <div class="oc-flex">
                <oc-icon v-if="dateRangeApplied" name="check" />
              </div>
            </oc-button>
          </li>
        </oc-list>
        <div
          class="date-filter-range-panel oc-py-s"
          :class="{ 'date-filter-range-panel-active': dateRangeClicked }"
        >
          <div class="oc-flex oc-flex-middle oc-flex-between oc-mb-m">
            <oc-button
              appearance="raw"
              class="date-filter-range-panel-back"
              :aria-label="$gettext('Go back to filter options')"
              @click="dateRangeClicked = false"
            >
              <oc-icon name="arrow-left-s" fill-type="line" />
            </oc-button>
            <span v-text="$gettext('Custom date range')" />
            <oc-button
              appearance="raw"
              class="date-filter-range-panel-close"
              :aria-label="$gettext('Close filter')"
              @click="filterChip.hideDrop()"
            >
              <oc-icon name="close" />
            </oc-button>
          </div>
          <div class="oc-mt-s">
            <oc-datepicker
              :label="$gettext('From')"
              type="date"
              :is-clearable="true"
              :current-date="fromDate"
              @date-changed="(value) => setDateRangeDate(value, 'from')"
            />
            <oc-datepicker
              :label="$gettext('To')"
              type="date"
              :is-clearable="true"
              :current-date="toDate"
              :min-date="fromDate ? fromDate : undefined"
              @date-changed="(value) => setDateRangeDate(value, 'to')"
            />
          </div>
          <div class="date-filter-apply-btn">
            <oc-button
              appearance="outline"
              variation="passive"
              size="small"
              :disabled="!dateRangeValid"
              @click="applyDateRangeFilter"
            >
              {{ $gettext('Apply') }}
            </oc-button>
          </div>
        </div>
      </template>
    </oc-filter-chip>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, unref } from 'vue'
import omit from 'lodash-es/omit'
import { useRoute, useRouteQuery, useRouter } from '../../composables'
import { formatDateFromDateTime } from '../../helpers'
import { queryItemAsString } from '../../composables/appDefaults'
import { DateTime } from 'luxon'
import { useGettext } from 'vue3-gettext'
import { type OcComponents } from '@ownclouders/design-system/helpers'

type Item = Record<string, string>

interface Props {
  filterLabel: string
  filterName: string
  items: Item[]
  idAttribute?: string
  displayNameAttribute?: string
}
interface Emits {
  (e: 'selectionChange', value: Item): void
}

const {
  filterLabel,
  filterName,
  items,
  idAttribute = 'id',
  displayNameAttribute = 'name'
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const router = useRouter()
const { current: currentLanguage } = useGettext()
const currentRoute = useRoute()
const selectedItem = ref<Item>()
const displayedItems = ref(items)
const fromDate = ref<DateTime>()
const toDate = ref<DateTime>()
const dateRangeClicked = ref(false)
const filterChip = ref<OcComponents['OcFilterChip']>()

const queryParam = `q_${filterName}`
const currentRouteQuery = useRouteQuery(queryParam)

const getId = (item: Item) => {
  return item[idAttribute as keyof Item]
}

const setRouteQuery = () => {
  return router.push({
    query: {
      ...omit(unref(currentRoute).query, [queryParam]),
      ...(unref(selectedItem) && {
        [queryParam]: getId(unref(selectedItem))
      })
    }
  })
}

const isItemSelected = (item: Item) => {
  return unref(selectedItem) && getId(unref(selectedItem)) === getId(item)
}

const resetDateRange = () => {
  fromDate.value = undefined
  toDate.value = undefined
}

const toggleItemSelection = async (item: Item) => {
  resetDateRange()
  if (isItemSelected(item)) {
    selectedItem.value = undefined
  } else {
    selectedItem.value = item
    unref(filterChip).hideDrop()
  }
  await setRouteQuery()
  emit('selectionChange', unref(selectedItem))
}

const clearFilter = () => {
  selectedItem.value = undefined
  dateRangeClicked.value = false
  resetDateRange()
  emit('selectionChange', unref(selectedItem))
  setRouteQuery()
}

const onShowDrop = () => {
  displayedItems.value = items
}

const onHideDrop = () => {
  dateRangeClicked.value = false
}

const setSelectedItemsBasedOnQuery = () => {
  const id = queryItemAsString(unref(currentRouteQuery))
  if (!id) {
    return
  }

  const selected = items.find((s) => getId(s) === id)
  if (selected) {
    selectedItem.value = selected
    return
  }

  if (unref(dateRangeApplied)) {
    const dateRange = id.replace('range:', '')
    const from = Number(dateRange.split(' - ')[0])
    const to = Number(dateRange.split(' - ')[1])
    fromDate.value = DateTime.fromMillis(from)
    toDate.value = DateTime.fromMillis(to)
    selectedItem.value = unref(dateRangeOption)
  }
}

const dateRangeApplied = computed(() =>
  queryItemAsString(unref(currentRouteQuery))?.startsWith('range:')
)

const dateRangeValid = computed(() => {
  if (!unref(fromDate) || !unref(toDate)) {
    return false
  }
  return unref(fromDate) <= unref(toDate)
})

const dateRangeOption = computed(() => {
  if (!unref(fromDate) || !unref(toDate)) {
    return null
  }

  const from = formatDateFromDateTime(unref(fromDate), currentLanguage, DateTime.DATE_SHORT)
  const to = formatDateFromDateTime(unref(toDate), currentLanguage, DateTime.DATE_SHORT)
  const fromDateMillis = unref(fromDate).toMillis()
  const toDateMillis = unref(toDate).toMillis()

  return {
    [idAttribute]: `range:${fromDateMillis} - ${toDateMillis}`,
    [displayNameAttribute]: `${from} - ${to}`
  }
})

const applyDateRangeFilter = async () => {
  selectedItem.value = unref(dateRangeOption)
  await setRouteQuery()
  unref(filterChip).hideDrop()
  emit('selectionChange', unref(selectedItem))
}

const setDateRangeDate = (
  { date, error }: { date: DateTime; error: Error },
  type: 'from' | 'to'
) => {
  if (error) {
    console.error(error)
    return
  }

  const prop = type === 'from' ? fromDate : toDate

  if (!date) {
    prop.value = undefined
    return
  }

  const formattedDate = type === 'from' ? date.startOf('day') : date.endOf('day')
  prop.value = formattedDate
}

defineExpose({ setSelectedItemsBasedOnQuery })

onMounted(() => {
  setSelectedItemsBasedOnQuery()
})
</script>

<style lang="scss">
.date-filter {
  overflow: hidden;

  &-list {
    li {
      &:first-child {
        margin-top: 0 !important;
      }
      &:last-child {
        margin-bottom: 0 !important;
      }
    }

    &-item {
      line-height: 1.5;
      gap: 8px;

      &:hover,
      &-active {
        background-color: var(--oc-color-background-hover) !important;
      }
    }

    &-hidden {
      min-height: 225px;
      visibility: hidden;
      transition: visibility 0.4s 0s;
    }
  }

  &-apply-btn {
    text-align: end;
  }

  &-range-panel {
    transform: translateX(100%);
    transition:
      transform 0.4s ease,
      visibility 0.4s 0s;
    visibility: hidden;
    position: absolute;
    width: calc(100% - var(--oc-space-medium));
    background: #fff;
    top: 0;
    color: var(--oc-color-swatch-passive-default);

    &-active {
      visibility: unset;
      transform: translateX(0);

      .oc-card {
        overflow: unset;
      }
    }
  }

  .oc-card {
    overflow: hidden;
  }

  .oc-date-picker label {
    font-size: var(--oc-font-size-small);
  }
}
</style>
