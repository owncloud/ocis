<template>
  <div class="oc-flex group-info oc-mb-l">
    <avatar-image class="oc-mb-m" :width="80" :userid="group.id" :user-name="group.displayName" />
    <span class="oc-text-muted group-info-display-name" v-text="group.displayName"></span>
    <span class="oc-text-muted" v-text="groupMembersText"></span>
  </div>
</template>
<script lang="ts" setup>
import { computed, unref } from 'vue'
import { Group } from '@ownclouders/web-client/graph/generated'
import { useGettext } from 'vue3-gettext'

interface Props {
  group: Group
}

const props = defineProps<Props>()
const group = computed<Group>(() => props.group)
const { $ngettext } = useGettext()
const groupMembersText = computed(() => {
  return $ngettext('%{groupCount} member', '%{groupCount} members', unref(group).members.length, {
    groupCount: unref(group).members.length.toString()
  })
})
</script>
<style lang="scss">
.gr-info {
  align-items: center;
  flex-direction: column;
}
.group-info-display-name {
  font-size: 1.5rem;
}
</style>
