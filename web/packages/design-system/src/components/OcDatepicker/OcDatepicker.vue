<template>
  <oc-text-input
    v-model="dateInputString"
    v-bind="$attrs"
    :label="label"
    type="date"
    :min="minDate?.toISODate()"
    :fix-message-line="true"
    :error-message="errorMessage"
    :clear-button-enabled="isClearable"
    :clear-button-accessible-label="$gettext('Clear date')"
    class="oc-date-picker"
  />
</template>

<script lang="ts" setup>
import { computed, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import { DateTime } from 'luxon'

interface Props {
  label: string
  isClearable?: boolean
  currentDate?: DateTime
  minDate?: DateTime
}

interface Emits {
  (e: 'dateChanged', value: { date: DateTime | null; error: boolean }): void
}
defineOptions({
  name: 'OcDatepicker',
  status: 'ready',
  release: '1.0.0'
})

const { label, isClearable = true, currentDate = null, minDate = null } = defineProps<Props>()
const emit = defineEmits<Emits>()
const { $gettext, current } = useGettext()
const dateInputString = ref<string>('')

const date = computed(() => {
  const date = DateTime.fromISO(unref(dateInputString)).endOf('day')
  return date.isValid ? date : null
})

const isMinDateUndercut = computed(() => {
  if (!minDate || !unref(date)) {
    return false
  }
  return unref(date) < minDate
})

const errorMessage = computed(() => {
  if (unref(isMinDateUndercut)) {
    return $gettext('The date must be after %{date}', {
      date: minDate.minus({ day: 1 }).setLocale(current).toLocaleString(DateTime.DATE_SHORT)
    })
  }
  return ''
})

watch(
  () => currentDate,
  () => {
    if (currentDate) {
      dateInputString.value = currentDate.toISODate()
      return
    }

    dateInputString.value = undefined
  },
  { immediate: true }
)

watch(
  date,
  () => {
    emit('dateChanged', { date: unref(date), error: unref(isMinDateUndercut) })
  },
  {
    deep: true
  }
)
</script>
<style lang="scss">
.oc-date-picker {
  input::-webkit-calendar-picker-indicator {
    cursor: pointer;
  }
}
</style>
<docs>
```js
<template>
  <div>
    <oc-datepicker :current-date="currentDate" :min-date="minDate" label="Enter or pick a date"
                   @date-changed="onDateChanged"/>
    <p v-if="selectedDate" v-text="selectedDate"/>
  </div>
</template>
<script>
  import {DateTime} from "luxon";

  export default {
    data: () => ({
      minDate: DateTime.now(), currentDate: DateTime.now(), selectedDate: ''
    }),
    methods: {
      onDateChanged({date}) {
        this.selectedDate = date ? date.toLocaleString(DateTime.DATE_FULL) : ''
      }
    }
  }
</script>
```
</docs>
