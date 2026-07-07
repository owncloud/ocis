<template>
  <oc-button
    v-if="availableLinkTypeOptions.length > 1"
    :id="`link-role-dropdown-toggle-${dropUuid}`"
    appearance="raw"
    gap-size="none"
    class="oc-text-left link-role-dropdown-toggle"
  >
    <span class="link-current-role" v-text="currentLinkRoleLabel || $gettext('Select a role')" />
    <oc-icon name="arrow-down-s" />
  </oc-button>
  <span
    v-else
    v-oc-tooltip="getLinkRoleByType(modelValue)?.description"
    class="link-current-role oc-mr-m"
    v-text="currentLinkRoleLabel"
  />
  <oc-drop
    v-if="availableLinkTypeOptions.length > 1"
    class="link-role-dropdown"
    :drop-id="`link-role-dropdown-${dropUuid}`"
    :toggle="`#link-role-dropdown-toggle-${dropUuid}`"
    padding-size="small"
    mode="click"
    :offset="dropOffset"
    close-on-click
  >
    <oc-list class="role-dropdown-list">
      <li v-for="(type, i) in availableLinkTypeOptions" :key="`role-dropdown-${i}`">
        <oc-button
          :id="`files-role-${getLinkRoleByType(type).id}`"
          :class="{
            selected: isSelectedType(type),
            'oc-background-primary-gradient': isSelectedType(type)
          }"
          :appearance="isSelectedType(type) ? 'raw-inverse' : 'raw'"
          :variation="isSelectedType(type) ? 'primary' : 'passive'"
          justify-content="space-between"
          class="oc-p-s"
          @click="updateSelectedType(type)"
        >
          <span class="oc-flex oc-flex-middle">
            <oc-icon
              :name="getLinkRoleByType(type).icon"
              class="oc-pl-s oc-pr-m"
              variation="inherit"
            />
            <span>
              <span
                class="role-dropdown-list-option-label oc-text-bold oc-display-block oc-width-1-1"
                v-text="$gettext(getLinkRoleByType(type).displayName)"
              />
              <span class="oc-text-small">{{ $gettext(getLinkRoleByType(type).description) }}</span>
            </span>
          </span>
          <span class="oc-flex">
            <oc-icon v-if="isSelectedType(type)" name="check" variation="inherit" />
          </span>
        </oc-button>
      </li>
    </oc-list>
  </oc-drop>
</template>

<script lang="ts" setup>
import { v4 as uuidV4 } from 'uuid'
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { useLinkTypes } from '../composables'

interface Emits {
  (event: 'update:modelValue', value: SharingLinkType): void
}

interface Props {
  modelValue: SharingLinkType
  availableLinkTypeOptions: SharingLinkType[]
  dropOffset?: string
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { $gettext } = useGettext()
const { getLinkRoleByType } = useLinkTypes()

const isSelectedType = (type: SharingLinkType) => {
  return props.modelValue === type
}

const updateSelectedType = (type: SharingLinkType) => {
  emit('update:modelValue', type)
}

const currentLinkRoleLabel = computed(() => {
  return $gettext(getLinkRoleByType(props.modelValue)?.displayName)
})

const dropUuid = uuidV4()
</script>

<style lang="scss" scoped>
@media (max-width: $oc-breakpoint-medium-default) {
  .link-role-dropdown {
    width: 100%;
  }
}

@media (min-width: $oc-breakpoint-medium-default) {
  .link-role-dropdown {
    width: 400px;
  }
}

.role-dropdown-list span {
  line-height: 1.3;
}

.role-dropdown-list li {
  margin: var(--oc-space-xsmall) 0;

  &:first-child {
    margin-top: 0;
  }

  &:last-child {
    margin-bottom: 0;
  }

  .oc-button {
    text-align: left;
    width: 100%;
    gap: var(--oc-space-medium);

    &:hover,
    &:focus {
      background-color: var(--oc-color-background-hover);
      text-decoration: none;
    }
  }

  .selected span {
    color: var(--oc-color-swatch-primary-contrast);
  }
}
</style>
