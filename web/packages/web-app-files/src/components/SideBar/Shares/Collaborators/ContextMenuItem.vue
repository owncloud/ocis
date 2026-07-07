<template>
  <div
    v-if="option.hasSwitch"
    class="action-menu-item item-has-switch oc-p-s oc-flex oc-flex-center"
  >
    <oc-icon :name="option.icon" fill-type="line" size="medium" variation="passive" />
    <oc-switch
      class="oc-ml-s oc-flex oc-width-1-1"
      :checked="option.isChecked.value"
      :class="option.class"
      :label="option.title"
      @update:checked="option.method"
    />
  </div>
  <oc-button
    v-else
    appearance="raw"
    class="oc-p-s action-menu-item"
    :class="option.class"
    :type="option.to ? 'router-link' : 'button'"
    :to="option.to"
    v-bind="option.additionalAttributes || {}"
    @click="option.method"
  >
    <oc-icon :name="option.icon" fill-type="line" size="medium" />
    <span v-text="option.title" />
  </oc-button>
</template>

<script lang="ts" setup>
import { EditOption } from './EditDropdown.vue'

interface Props {
  option: EditOption
}
const { option } = defineProps<Props>()
</script>
