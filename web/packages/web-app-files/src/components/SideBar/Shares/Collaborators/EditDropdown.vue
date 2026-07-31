<template>
  <div class="oc-flex oc-flex-middle">
    <oc-button
      :id="editShareBtnId"
      v-oc-tooltip="dropButtonTooltip"
      class="collaborator-edit-dropdown-options-btn"
      :aria-label="
        isLocked ? dropButtonTooltip : $gettext('Open context menu with share editing options')
      "
      appearance="raw"
      :disabled="isLocked"
    >
      <oc-icon name="more-2" :accessible-label="$gettext('Access details')" />
    </oc-button>
    <oc-drop
      ref="expirationDateDrop"
      :toggle="'#' + editShareBtnId"
      mode="click"
      padding-size="small"
    >
      <oc-list class="collaborator-edit-dropdown-options-list" :aria-label="shareEditOptions">
        <li v-for="(option, i) in options" :key="i" class="oc-rounded oc-menu-item-hover">
          <context-menu-item :option="option" />
        </li>
        <li v-if="sharedParentRoute" class="oc-rounded oc-menu-item-hover">
          <context-menu-item :option="navigateToParentOption" />
        </li>
      </oc-list>
      <oc-list
        v-if="canEdit"
        class="collaborator-edit-dropdown-options-list collaborator-edit-dropdown-options-list-remove"
      >
        <li
          class="oc-rounded oc-menu-item-hover"
          :class="{ 'oc-pt-s': options.length > 0 || sharedParentRoute }"
        >
          <context-menu-item :option="removeShareOption" />
        </li>
      </oc-list>
    </oc-drop>
    <oc-info-drop
      ref="accessDetailsDrop"
      class="share-access-details-drop"
      v-bind="{
        title: $gettext('Access details'),
        list: accessDetails
      }"
      mode="manual"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, inject, Ref, unref, useTemplateRef } from 'vue'
import { DateTime } from 'luxon'
import { ContextualHelperDataListItem, uniqueId } from '@ownclouders/design-system/helpers'
import { OcDrop, OcInfoDrop } from '@ownclouders/design-system/components'
import { Resource } from '@ownclouders/web-client'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import { useConfigStore, useModals } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import DatePickerModal from '../../../Modals/DatePickerModal.vue'
import { RouteLocationNamedRaw } from 'vue-router'
import ContextMenuItem from './ContextMenuItem.vue'

export type EditOption = {
  icon: string
  title: string
  additionalAttributes?: Record<string, string>
  class?: string
  hasSwitch?: boolean
  isChecked?: Ref<boolean>
  method?: (args?: unknown) => void
  to?: RouteLocationNamedRaw
}

interface Props {
  expirationDate?: string
  shareCategory?: 'user' | 'group' | null
  canEdit: boolean
  accessDetails: ContextualHelperDataListItem[]
  isShareDenied?: boolean
  deniable?: boolean
  isLocked?: boolean
  sharedParentRoute?: RouteLocationNamedRaw
}
interface Emits {
  (e: 'expirationDateChanged', payload: { expirationDateTime: DateTime | null }): void
  (e: 'removeShare'): void
  (e: 'setDenyShare', value: boolean): void
  (e: 'notifyShare'): void
}
const {
  expirationDate = undefined,
  shareCategory = null,
  canEdit,
  accessDetails,
  isShareDenied = false,
  deniable = false,
  isLocked = false,
  sharedParentRoute = undefined
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const language = useGettext()
const { $gettext } = language
const configStore = useConfigStore()
const { dispatchModal } = useModals()
const expirationDateDrop = useTemplateRef<typeof OcDrop>('expirationDateDrop')
const accessDetailsDrop = useTemplateRef<typeof OcInfoDrop>('accessDetailsDrop')

const resource = inject<Ref<Resource>>('resource')

const toggleShareDenied = (value: boolean) => {
  emit('setDenyShare', value)
}

function removeExpirationDate() {
  emit('expirationDateChanged', { expirationDateTime: null })
  unref(expirationDateDrop).hide()
}
function showDatePickerModal() {
  const currentDate = DateTime.fromISO(expirationDate)

  dispatchModal({
    title: $gettext('Set expiration date'),
    hideActions: true,
    customComponent: DatePickerModal,
    customComponentAttrs: () => ({
      currentDate: currentDate.isValid ? currentDate : null,
      minDate: DateTime.now()
    }),
    onConfirm: (expirationDateTime: DateTime) => {
      emit('expirationDateChanged', {
        expirationDateTime
      })
    }
  })
}
const dropButtonTooltip = computed(() => {
  if (isLocked) {
    return $gettext('Resource is temporarily locked, unable to manage share')
  }

  return ''
})

const navigateToParentOption = computed<EditOption>(() => {
  return {
    title: $gettext('Navigate to parent'),
    icon: 'folder-shared',
    class: 'navigate-to-parent',
    to: sharedParentRoute
  }
})

const removeShareOption = computed<EditOption>(() => {
  return {
    title: isProjectSpaceResource(unref(resource))
      ? $gettext('Remove member')
      : $gettext('Remove share'),
    method: () => {
      emit('removeShare')
    },
    class: 'remove-share',
    icon: 'delete-bin-5',
    additionalAttributes: {
      'data-testid': 'collaborator-remove-share-btn'
    }
  }
})

const options = computed<EditOption[]>(() => {
  const result: EditOption[] = [
    {
      title: $gettext('Access details'),
      method: () => accessDetailsDrop.value.$refs.drop.show(),
      icon: 'information',
      class: 'show-access-details'
    }
  ]

  if (deniable) {
    result.push({
      title: $gettext('Deny access'),
      method: toggleShareDenied,
      icon: 'stop-circle',
      class: 'deny-share',
      hasSwitch: true,
      isChecked: computed(() => isShareDenied)
    })
  }

  if (canEdit && unref(isExpirationSupported)) {
    result.push({
      title: unref(isExpirationDateSet)
        ? $gettext('Edit expiration date')
        : $gettext('Set expiration date'),
      class: 'set-expiration-date recipient-datepicker-btn',
      icon: 'calendar-event',
      method: showDatePickerModal
    })
  }

  if (unref(isRemoveExpirationPossible)) {
    result.push({
      title: $gettext('Remove expiration date'),
      class: 'remove-expiration-date',
      icon: 'calendar-close',
      method: removeExpirationDate
    })
  }

  if (configStore.options.isRunningOnEos) {
    result.push({
      title: $gettext('Notify via mail'),
      method: () => emit('notifyShare'),
      icon: 'mail',
      class: 'notify-via-mail'
    })
  }

  return result
})

const editShareBtnId = computed(() => {
  return uniqueId('files-collaborators-edit-button-')
})
const shareEditOptions = computed(() => {
  return $gettext('Context menu of the share')
})

const editingUser = computed(() => {
  return shareCategory === 'user'
})

const editingGroup = computed(() => {
  return shareCategory === 'group'
})

const isExpirationSupported = computed(() => {
  return unref(editingUser) || unref(editingGroup)
})

const isExpirationDateSet = computed(() => {
  return !!expirationDate
})

const isRemoveExpirationPossible = computed(() => {
  return canEdit && unref(isExpirationSupported) && unref(isExpirationDateSet)
})
</script>
<style lang="scss">
.collaborator-edit-dropdown-options-list {
  &-remove {
    margin-top: var(--oc-space-small) !important;
    border-top: 1px solid var(--oc-color-border) !important;
  }

  .action-menu-item {
    width: 100%;
    justify-content: flex-start;
    color: var(--oc-color-swatch-passive-default);
    gap: var(--oc-space-small);
  }
}
.share-access-details-drop {
  dl {
    display: grid;
    grid-template-columns: max-content auto;
    column-gap: var(--oc-space-medium);
    row-gap: var(--oc-space-xsmall);
  }
  dt {
    grid-column-start: 1;
  }
  dd {
    grid-column-start: 2;
    margin-left: var(--oc-space-medium);
  }
}
</style>
