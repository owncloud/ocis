<template>
  <div class="oc-flex user-info oc-mb-l">
    <avatar-image class="oc-mb-m" :width="80" :userid="user.id" :user-name="user.displayName" />
    <span v-text="user.onPremisesSamAccountName"></span>
    <span class="oc-text-muted user-info-display-name" v-text="user.displayName"></span>
  </div>
</template>
<script lang="ts" setup>
import { User } from '@ownclouders/web-client/graph/generated'

interface Props {
  user: User
}
const { user } = defineProps<Props>()
</script>
<style lang="scss">
.user-info {
  align-items: center;
  flex-direction: column;
}
.user-info-display-name {
  font-size: 1.5rem;
}
</style>
