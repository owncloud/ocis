<template>
  <span :id="`files-role-${(role as ShareRole).id}`" class="roles-select-role-item">
    <span
      class="oc-text-bold oc-display-block oc-width-1-1"
      v-text="$gettext((role as ShareRole).displayName)"
    />
    <span
      class="oc-m-rm oc-text-small oc-display-block"
      v-text="$gettext((role as ShareRole).description)"
    />
  </span>
</template>

<script lang="ts" setup>
import { ShareRole } from '@ownclouders/web-client'

interface Props {
  role: ShareRole
}
const { role } = defineProps<Props>()
</script>

<style lang="scss" scoped>
.roles-select-role-item {
  text-align: left;

  span {
    line-height: 1.3;
  }
}
</style>
