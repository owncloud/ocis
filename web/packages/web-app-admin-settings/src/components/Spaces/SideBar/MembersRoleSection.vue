<template>
  <ul class="oc-list">
    <li
      v-for="(m, index) in props.members"
      :key="index"
      class="oc-flex oc-flex-middle oc-mb-s"
      data-testid="space-members-list"
    >
      <oc-avatar
        v-if="m.grantedTo.user"
        :user-name="m.grantedTo.user.displayName"
        :width="36"
        class="oc-mr-s"
      /><oc-avatar-item
        v-else
        :width="36"
        icon-size="medium"
        :icon="groupIcon"
        name="group"
        class="oc-mr-s"
      />
      {{ (m.grantedTo.user || m.grantedTo.group).displayName }}
    </li>
  </ul>
</template>
<script lang="ts" setup>
import { computed } from 'vue'
import { ShareTypes, SpaceMember } from '@ownclouders/web-client'

interface Props {
  members: SpaceMember[]
}

const props = defineProps<Props>()
const groupIcon = computed(() => ShareTypes.group.icon)
</script>
