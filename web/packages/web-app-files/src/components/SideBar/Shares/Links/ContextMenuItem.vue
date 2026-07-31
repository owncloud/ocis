<template>
  <oc-button
    appearance="raw"
    class="oc-p-s action-menu-item"
    :type="option.to ? 'router-link' : 'button'"
    :to="option.to"
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
