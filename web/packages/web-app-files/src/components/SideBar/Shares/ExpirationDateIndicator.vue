<template>
  <div class="oc-flex oc-flex-center expiration-date-indicator">
    <oc-icon
      v-oc-tooltip="expirationDateTooltip"
      :accessible-label="screenreaderShareExpiration"
      name="calendar-event"
      fill-type="line"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { DateTime } from 'luxon'
import { formatDateFromDateTime, formatRelativeDateFromDateTime } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

interface Props {
  expirationDate?: DateTime
}
const { expirationDate = null } = defineProps<Props>()
const { $gettext, current: currentLanguage } = useGettext()

const expirationDateRelative = computed(() => {
  return formatRelativeDateFromDateTime(expirationDate, currentLanguage)
})

const dateExpire = computed(() => {
  return formatDateFromDateTime(expirationDate, currentLanguage)
})

const expirationDateTooltip = computed(() => {
  return $gettext(
    'Expires %{timeToExpiry} (%{expiryDate})',
    { timeToExpiry: unref(expirationDateRelative), expiryDate: unref(dateExpire) },
    true
  )
})

const screenreaderShareExpiration = computed(() => {
  return $gettext('Share expires %{ expiryDateRelative } (%{ expiryDate })', {
    expiryDateRelative: unref(expirationDateRelative),
    expiryDate: unref(dateExpire)
  })
})
</script>
