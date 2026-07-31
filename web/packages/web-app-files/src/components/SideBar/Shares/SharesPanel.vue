<template>
  <div>
    <oc-loader v-if="sharesLoading" :aria-label="$gettext('Loading list of shares')" />
    <div v-else-if="hasSharesLoadingFailed" class="oc-text-center oc-pt-xl">
      <oc-icon name="group" variation="danger" size="xxlarge" />
      <p class="oc-text-danger">
        {{ errorMessage }}
      </p>
    </div>
    <template v-else>
      <space-members v-if="showSpaceMembers" class="oc-background-muted oc-p-m oc-mb-s" />
      <file-shares v-else class="oc-background-muted oc-p-m oc-mb-s" />
      <file-links v-if="showLinks" class="oc-background-muted oc-p-m" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, unref } from 'vue'
import FileLinks from './FileLinks.vue'
import FileShares from './FileShares.vue'
import SpaceMembers from './SpaceMembers.vue'
import { useSharesStore } from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'
import { useGettext } from 'vue3-gettext'

defineOptions({ name: 'SharesPanel' })

const { showSpaceMembers = false, showLinks = false } = defineProps<{
  showSpaceMembers?: boolean
  showLinks?: boolean
}>()

const { $gettext } = useGettext()

const sharesStore = useSharesStore()
const { loading: sharesLoading, hasLoadingFailed: hasSharesLoadingFailed } =
  storeToRefs(sharesStore)

const errorMessage = computed(() => {
  if (!unref(hasSharesLoadingFailed)) {
    return ''
  }

  if (unref(showSpaceMembers)) {
    return $gettext('Loading members failed. Try again later.')
  }

  return $gettext('Loading shares failed. Try again later.')
})
</script>
