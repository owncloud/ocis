<template>
  <div class="quota-select-batch-action-form">
    <oc-select
      ref="select"
      :model-value="selectedOption"
      :selectable="optionSelectable"
      taggable
      push-tags
      :clearable="false"
      :options="options"
      :create-option="createOption"
      option-label="displayValue"
      :label="$gettext('Quota')"
      v-bind="$attrs"
      @update:model-value="onUpdate"
    >
      <template #selected-option="{ displayValue }">
        <oc-icon v-if="$attrs['read-only']" name="lock" class="oc-mr-xs" size="small" />
        <span v-text="displayValue" />
      </template>
      <template #search="{ attributes, events }">
        <input class="vs__search" v-bind="attributes" v-on="events" />
      </template>
      <template #option="{ displayValue, error }">
        <div class="oc-flex oc-flex-between">
          <span v-text="displayValue" />
        </div>
        <div v-if="error" class="oc-text-input-danger">{{ error }}</div>
      </template>
    </oc-select>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useGettext } from 'vue3-gettext'
import { isNumber } from 'lodash-es'
import { formatFileSize } from '../helpers'

type Option = {
  value: number
  displayValue: string
  selectable?: boolean
}

interface Props {
  totalQuota?: number
  maxQuota?: number
}

interface Emits {
  (e: 'selectedOptionChange', option: Option): void
}

const props = withDefaults(defineProps<Props>(), {
  totalQuota: 0,
  maxQuota: 0
})
const emit = defineEmits<Emits>()
const selectedOption = ref<Option>(undefined)
const options = ref<Option[]>([])
const { $gettext, current } = useGettext()
function onUpdate(event: Option) {
  selectedOption.value = event
  emit('selectedOptionChange', selectedOption.value)
}
function optionSelectable(option: Option) {
  return option.selectable !== false
}
function isValueValidNumber(value: string | number) {
  if (isNumber(value)) {
    return value > 0
  }

  const optionIsNumberRegex = /^[0-9]\d*(([.,])\d+)?$/g
  return optionIsNumberRegex.test(value)
}
function createOption(option: string) {
  option = option.replace(',', '.')

  if (!isValueValidNumber(option)) {
    return {
      displayValue: option,
      value: option,
      error: $gettext('Please enter only numbers'),
      selectable: false
    }
  }
  const value = parseFloat(option) * Math.pow(10, 9)

  if (value > quotaLimit.value) {
    return {
      value,
      displayValue: getFormattedFileSize(value),
      error: $gettext('Please enter a value equal to or less than %{ quotaLimit }', {
        quotaLimit: getFormattedFileSize(quotaLimit.value).toString()
      }),

      selectable: false
    }
  }

  return {
    value,
    displayValue: getFormattedFileSize(value)
  }
}
function setOptions() {
  let availableOptions = [...DEFAULT_OPTIONS.value]

  if (props.maxQuota) {
    availableOptions = availableOptions.filter((availableOption) => {
      if (props.totalQuota === 0 && availableOption.value === 0) {
        availableOption.selectable = false
        return true
      }
      return availableOption.value !== 0 && availableOption.value <= props.maxQuota
    })
  }

  const selectedQuotaInOptions = availableOptions.find(
    (option) => option.value === props.totalQuota
  )

  if (!selectedQuotaInOptions) {
    availableOptions.push({
      displayValue: getFormattedFileSize(props.totalQuota),
      value: props.totalQuota,
      selectable: props.totalQuota <= quotaLimit.value
    })
  }

  // Sort options and make sure that unlimited is at the end
  availableOptions = [
    ...availableOptions.filter((o) => o.value).sort((a, b) => a.value - b.value),
    ...availableOptions.filter((o) => !o.value)
  ]
  options.value = availableOptions
}
function getFormattedFileSize(value: number) {
  const formattedFilesize = formatFileSize(value, current)
  return !isValueValidNumber(value) ? value.toString() : formattedFilesize
}
watch(
  () => props.totalQuota,
  () => {
    const option = options.value.find((o) => o.value === props.totalQuota)
    if (option) {
      selectedOption.value = option
    }
  }
)
const quotaLimit = computed(() => {
  return props.maxQuota || 1e15
})
const DEFAULT_OPTIONS = computed<Option[]>(() => {
  return [
    {
      value: Math.pow(10, 9),
      displayValue: getFormattedFileSize(Math.pow(10, 9))
    },
    {
      value: 2 * Math.pow(10, 9),
      displayValue: getFormattedFileSize(2 * Math.pow(10, 9))
    },
    {
      value: 5 * Math.pow(10, 9),
      displayValue: getFormattedFileSize(5 * Math.pow(10, 9))
    },
    {
      value: 10 * Math.pow(10, 9),
      displayValue: getFormattedFileSize(10 * Math.pow(10, 9))
    },
    {
      value: 50 * Math.pow(10, 9),
      displayValue: getFormattedFileSize(50 * Math.pow(10, 9))
    },
    {
      value: 100 * Math.pow(10, 9),
      displayValue: getFormattedFileSize(100 * Math.pow(10, 9))
    },
    {
      displayValue: $gettext('No restriction'),
      value: 0
    }
  ]
})
onMounted(() => {
  setOptions()
  selectedOption.value = options.value.find((o) => o.value === props.totalQuota)
})
</script>
