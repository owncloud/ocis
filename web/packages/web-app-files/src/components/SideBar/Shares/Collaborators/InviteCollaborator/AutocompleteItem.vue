<template>
  <div
    :data-testid="`recipient-autocomplete-item-${item.displayName}`"
    class="oc-flex oc-flex-middle oc-py-xs"
    :class="collaboratorClass"
  >
    <avatar-image
      v-if="isAnyUserShareType"
      class="oc-mr-s"
      :width="36"
      :userid="item.id"
      :user-name="item.displayName"
    />
    <oc-avatar-item
      v-else
      :width="36"
      :name="shareTypeKey"
      :icon="shareTypeIcon"
      icon-size="large"
      icon-color="var(--oc-color-text)"
      background="transparent"
      class="oc-mr-s"
    />
    <div class="files-collaborators-autocomplete-user-text oc-text-truncate">
      <span class="files-collaborators-autocomplete-username" v-text="item.displayName" />
      <template v-if="!isAnyPrimaryShareType">
        <span
          class="files-collaborators-autocomplete-share-type"
          v-text="`(${$gettext(shareType.label)})`"
        />
      </template>
      <div
        v-if="additionalInfo"
        class="files-collaborators-autocomplete-additionalInfo"
        v-text="`${additionalInfo}`"
      />
      <div
        v-if="externalIssuer"
        class="files-collaborators-autocomplete-externalIssuer"
        v-text="`${externalIssuer}`"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { CollaboratorAutoCompleteItem, ShareTypes } from '@ownclouders/web-client'

interface Props {
  item: CollaboratorAutoCompleteItem
}

const props = defineProps<Props>()

const additionalInfo = computed(() => {
  return (
    props.item.attributes?.join(' · ') || props.item.mail || props.item.onPremisesSamAccountName
  )
})

const externalIssuer = computed(() => {
  if (props.item.shareType === ShareTypes.remote.value) {
    return props.item.identities?.[0]?.issuer
  }
  return ''
})
const shareType = computed(() => {
  return ShareTypes.getByValue(props.item.shareType)
})

const shareTypeIcon = computed(() => {
  return unref(shareType).icon
})

const shareTypeKey = computed(() => {
  return unref(shareType).key
})

const isAnyUserShareType = computed(() => {
  return ShareTypes.user.key === unref(shareType).key
})

const isAnyPrimaryShareType = computed(() => {
  return [ShareTypes.user.key, ShareTypes.group.key].includes(unref(shareType).key)
})

const collaboratorClass = computed(() => {
  return `files-collaborators-search-${unref(shareType).key}`
})
</script>

<style lang="scss">
.files-collaborators-autocomplete-additionalInfo,
.files-collaborators-autocomplete-externalIssuer {
  font-size: var(--oc-font-size-small);
  white-space: normal;
}
</style>
