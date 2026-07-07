<template>
  <div class="oc-status-indicators">
    <template v-for="indicator in indicators">
      <oc-button
        v-if="hasHandler(indicator) && !disableHandler"
        :id="indicator.id"
        :key="`${indicator.id}-handler`"
        v-oc-tooltip="$gettext(indicator.label)"
        class="oc-status-indicators-indicator oc-ml-xs"
        :aria-label="$gettext(indicator.label)"
        :aria-describedby="getIndicatorDescriptionId(indicator)"
        appearance="raw"
        :data-testid="indicator.id"
        :data-test-indicator-type="indicator.type"
        @click="indicator.handler(resource)"
      >
        <oc-icon
          :name="indicator.icon"
          size="small"
          :fill-type="indicator.fillType"
          variation="inherit"
        />
      </oc-button>
      <oc-icon
        v-else
        :id="indicator.id"
        :key="indicator.id"
        v-oc-tooltip="$gettext(indicator.label)"
        tabindex="-1"
        size="small"
        class="oc-status-indicators-indicator oc-ml-xs"
        :name="indicator.icon"
        :fill-type="indicator.fillType"
        :accessible-label="$gettext(indicator.label)"
        :aria-describedby="getIndicatorDescriptionId(indicator)"
        :data-testid="indicator.id"
        :data-test-indicator-type="indicator.type"
      />
      <p
        v-if="getIndicatorDescriptionId(indicator)"
        :id="getIndicatorDescriptionId(indicator)"
        :key="getIndicatorDescriptionId(indicator)"
        class="oc-invisible-sr"
        v-text="$gettext(indicator.accessibleDescription)"
      />
    </template>
  </div>
</template>

<script lang="ts" setup>
import { ref, unref } from 'vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcButton from '../OcButton/OcButton.vue'
import { uniqueId } from '../../helpers'

/**
 * OcStatusIndicators - Status indicators which can be attached to a resource.
 *
 * @component
 * @example
 * ```vue
 * <oc-status-indicators :resource="resource" :indicators="indicators" />
 * ```
 *
 * @prop {Object} resource - A resource to which the indicators are attached.
 * @prop {Array.<Indicator>} indicators - An array of indicators to be displayed.
 * @prop {boolean} [disableHandler=false] - Disables the handler for all indicators. Useful for disabled resources.
 *
 * @typedef {Object} Indicator
 * @property {string} id - Id of the indicator.
 * @property {string} icon - Icon of the indicator.
 * @property {string} label - String to be used as an accessible label and tooltip for the indicator.
 * @property {Function} [handler] - An action to be triggered when the indicator is clicked. Receives the resource.
 * @property {string} [accessibleDescription] - A string to be used as an accessible description for the indicator. It renders an element only visible for screen readers to provide additional context.
 * @property {boolean} [visible] - Visibility of the indicator.
 * @property {string} [type] - Type of the indicator.
 * @property {string} [fillType] - Fill type of the indicator.
 *
 * @method hasHandler
 * @param {Indicator} indicator - The indicator to check for a handler.
 * @returns {boolean} - Returns true if the indicator has a handler.
 *
 * @method getIndicatorDescriptionId
 * @param {Indicator} indicator - The indicator to get the description ID for.
 * @returns {string|null} - Returns the description ID if available, otherwise null.
 */

type Indicator = {
  id: string
  icon: string
  label: string
  handler?: any
  accessibleDescription?: string
  visible?: boolean
  type?: string
  fillType?: string
}

interface Props {
  resource: Record<string, any>
  indicators: Indicator[]
  disableHandler?: boolean
}

defineOptions({
  name: 'OcStatusIndicators',
  status: 'ready',
  release: '2.0.1'
})

const { resource, indicators, disableHandler = false } = defineProps<Props>()
const accessibleDescriptionIds = ref({} as Record<string, string>)

const hasHandler = (indicator: Indicator): boolean => {
  return Object.prototype.hasOwnProperty.call(indicator, 'handler')
}

const getIndicatorDescriptionId = (indicator: Indicator): string | null => {
  if (!indicator.accessibleDescription) {
    return null
  }

  if (!unref(accessibleDescriptionIds)[indicator.id]) {
    unref(accessibleDescriptionIds)[indicator.id] = uniqueId('oc-indicator-description-')
  }

  return unref(accessibleDescriptionIds)[indicator.id]
}
</script>

<style lang="scss">
.oc-status-indicators {
  align-items: center;
  display: flex;
  justify-content: flex-end;
}
</style>

<docs>
```js
<template>
  <oc-status-indicators :resource="resource" :indicators="indicators" />
</template>
<script>
  export default {
    data: () => ({
      resource: {
        name: "Documents",
        path: "/"
      },
      indicators: [
        {
          id: 'files-sharing',
          label: "Shared with other people",
          icon: 'group',
          handler: (resource, indicatorId) => alert(`Resource: ${resource.name}, indicator: ${indicatorId}`)
        },
        {
          id: 'file-link',
          label: "Shared via link",
          icon: 'links',
          handler: (resource, indicatorId) => alert(`Resource: ${resource.name}, indicator: ${indicatorId}`)
        }
      ]
    }),
  }
</script>
```
</docs>
