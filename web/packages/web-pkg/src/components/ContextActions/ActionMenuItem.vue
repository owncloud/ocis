<template>
  <li v-oc-tooltip="componentProps.disabled ? action.disabledTooltip?.(actionOptions) : ''">
    <oc-button
      v-oc-tooltip="showTooltip || action.hideLabel ? action.label(actionOptions) : ''"
      :type="componentType"
      v-bind="componentProps"
      :class="[action.class, 'action-menu-item', 'oc-py-s', 'oc-px-m', 'oc-width-1-1']"
      :aria-label="ariaLabel"
      data-testid="action-handler"
      :size="size"
      justify-content="left"
      :title="action.label(actionOptions)"
      v-on="componentListeners"
    >
      <oc-img
        v-if="action.img"
        data-testid="action-img"
        :src="action.img"
        alt=""
        class="oc-icon oc-icon-m"
      />
      <oc-img
        v-else-if="hasExternalImageIcon"
        data-testid="action-img"
        :src="action.icon"
        alt=""
        class="oc-icon oc-icon-m"
      />
      <oc-icon
        v-else-if="action.icon"
        data-testid="action-icon"
        :name="action.icon"
        :fill-type="action.iconFillType || 'line'"
        :size="size"
      />
      <span
        v-if="!action.hideLabel"
        class="oc-files-context-action-label oc-flex"
        data-testid="action-label"
      >
        <span v-text="action.label(actionOptions)" />
        <span
          v-if="action.showOpenInNewTabHint"
          class="oc-text-muted oc-text-xsmall"
          v-text="openInNewTabHint"
        />
      </span>
      <span
        v-if="action.shortcut && shortcutHint"
        class="oc-files-context-action-shortcut"
        v-text="action.shortcut"
      />
    </oc-button>
  </li>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { Action, ActionOptions, useConfigStore } from '../../composables'
import { useGettext } from 'vue3-gettext'
import { storeToRefs } from 'pinia'

interface Props {
  action: Action
  actionOptions: ActionOptions
  size?: string
  appearance?: string
  variation?: string
  shortcutHint?: boolean
  showTooltip?: boolean
  hasLimitedScreenSpace?: boolean
}
const {
  action,
  actionOptions,
  size = 'medium',
  appearance = 'raw',
  variation = 'passive',
  shortcutHint = true,
  showTooltip = false,
  hasLimitedScreenSpace = false
} = defineProps<Props>()
const { $gettext } = useGettext()
const configStore = useConfigStore()
const { options } = storeToRefs(configStore)

const componentType = computed<string>(() => {
  if (Object.hasOwn(action, 'route')) {
    return 'router-link'
  }
  if (Object.hasOwn(action, 'href')) {
    return 'a'
  }
  if (Object.hasOwn(action, 'handler')) {
    return 'button'
  }
  console.warn('ActionMenuItem: No handler, route or href callback found in action', action)
  return 'button'
})

const ariaLabel = computed<string | null>(() => {
  if (componentProps.value.disabled && action.disabledTooltip) {
    return action.disabledTooltip(actionOptions)
  }

  if (hasLimitedScreenSpace) {
    return action.label(actionOptions)
  }

  return ''
})
const componentProps = computed(() => {
  const properties = {
    appearance: action.appearance || appearance,
    variation: action.variation || variation,
    ...(action.isDisabled && {
      disabled: action.isDisabled(actionOptions)
    }),
    ...(action.id && { id: action.id })
  }

  return {
    ...properties,
    ...(unref(componentType) === 'router-link' && {
      to: action.route(actionOptions)
    }),
    ...(unref(componentType) === 'a' && {
      href: action.href(actionOptions)
    }),
    ...(['router-link', 'a'].includes(unref(componentType)) && {
      target: options.value.cernFeatures ? '_blank' : '_self'
    })
  }
})

const isMacOs = computed(() => {
  return window.navigator.userAgent.includes('Mac')
})

const openInNewTabHint = computed(() => {
  return $gettext(
    'Hold %{key} and click to open in new tab',
    { key: unref(isMacOs) ? '⌘' : $gettext('ctrl') },
    true
  )
})

const hasExternalImageIcon = computed(() => {
  return action.icon && /^https?:\/\//i.test(action.icon)
})
const componentListeners = computed(() => {
  if (typeof action.handler !== 'function') {
    return {}
  }

  const callback = () =>
    action.handler({
      ...actionOptions
    })
  if (action.keepOpen) {
    return {
      click: (event: Event) => {
        event.stopPropagation()
        callback()
      }
    }
  }
  return {
    click: callback
  }
})
</script>
<style lang="scss">
.action-menu-item {
  vertical-align: middle;
}

.oc-files-context-action-label {
  flex-direction: column;
}

.oc-files-context-action-shortcut {
  justify-content: right !important;
  font-size: var(--oc-font-size-small);
}
</style>
