<template>
  <oc-datepicker
    :label="$gettext('Expiration date')"
    type="date"
    :min-date="minDate"
    :current-date="currentDate"
    :is-clearable="isClearable"
    @date-changed="onDateChanged"
  />

  <div class="link-modal-actions oc-flex oc-flex-right oc-flex-middle oc-mt-s">
    <oc-button
      class="oc-modal-body-actions-cancel oc-ml-s"
      appearance="outline"
      variation="passive"
      @click="$emit('cancel')"
      >{{ $gettext('Cancel') }}
    </oc-button>
    <oc-button
      :disabled="confirmDisabled"
      class="oc-modal-body-actions-confirm oc-ml-s"
      appearance="filled"
      variation="primary"
      @click="$emit('confirm', dateTime)"
      >{{ $gettext('Confirm') }}
    </oc-button>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { DateTime } from 'luxon'

interface Props {
  currentDate?: DateTime
  minDate?: DateTime
  isClearable?: boolean
}
interface Emits {
  (e: 'confirm', value: DateTime): void
  (e: 'cancel'): void
}
defineEmits<Emits>()
const { currentDate = null, minDate = null, isClearable = true } = defineProps<Props>()
const dateTime = ref<DateTime>()
const confirmDisabled = ref(true)
const onDateChanged = ({ date, error }: { date: DateTime; error: boolean }) => {
  confirmDisabled.value = error || !date
  dateTime.value = date
}
</script>
