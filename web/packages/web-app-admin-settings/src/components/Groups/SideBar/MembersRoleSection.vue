<template>
  <ul class="oc-list">
    <li
      v-for="(member, index) in groupMembers"
      :key="index"
      class="oc-flex oc-flex-middle oc-mb-s"
      data-testid="group-members-list"
    >
      <oc-avatar :user-name="member.displayName" :width="36" class="oc-mr-s" />
      <span>{{ member.displayName }}</span>
    </li>
  </ul>
</template>
<script lang="ts" setup>
import { User } from '@ownclouders/web-client/graph/generated'

interface Props {
  groupMembers?: User[]
}

const { groupMembers = [] } = defineProps<Props>()
</script>
