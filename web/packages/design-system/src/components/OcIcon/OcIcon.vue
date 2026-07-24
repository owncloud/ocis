<template>
  <component
    :is="type"
    :class="[
      { 'oc-button-reset': type === 'button' },
      'oc-icon',
      sizeClass(size),
      variationClass(variation)
    ]"
  >
    <inline-svg
      :src="nameWithFillType"
      :transform-source="transformSvgElement"
      :role="accessibleLabel !== '' ? 'img' : null"
      :aria-hidden="accessibleLabel === '' ? 'true' : null"
      :aria-labelledby="accessibleLabel === '' ? null : svgTitleId"
      :focusable="accessibleLabel === '' ? 'false' : null"
      :style="color !== '' ? { fill: color } : {}"
      @loaded="$emit('loaded')"
    />
  </component>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import InlineSvg from 'vue-inline-svg'
import { type AvailableSizeType } from '../../helpers'
import { getSizeClass, uniqueId } from '../../helpers'

/**
 * @component OcIcon
 * @description
 * Icons are used to visually communicate core parts of the product and available actions.
 * They can act as wayfinding tools to help users more easily understand where they are in the product.
 *
 * @Accessibility
 * - Pass a label to the icon via the `accessibleLabel` property. The component will automatically add a `title` element which is also referenced by its ID via `aria-labelledby`.
 * - Omit `accessibleLabel` if the icon is decorative. In this case:
 *   - `aria-hidden` is set to `true`.
 *   - `focusable` is set to `false`.
 *   - All aria-related properties are removed or emptied.
 *
 * @props
 * @prop {string} [name='info'] - The name of the icon to display.
 * @prop {'fill'|'line'|'none'} [fillType='fill'] - The fill type of the icon.
 * @prop {string} [accessibleLabel=''] - Descriptive text for screen readers. Leave empty for decorative icons.
 * @prop {string} [type='span'] - The HTML element used for the icon.
 * @prop {AvailableSizeType} [size='medium'] - The size of the icon. Options: `xsmall`, `small`, `medium`, `large`, `xlarge`, `xxlarge`.
 * @prop {'passive'|'primary'|'danger'|'success'|'warning'|'brand'|'inherit'} [variation='passive'] - Style variation for additional meaning.
 * @prop {string} [color=''] - Overwrites the color of the icon.
 *
 * @emits
 * @event loaded - Emitted when the SVG is successfully loaded.
 *
 * @example
 * <template>
 *   <!-- Example of a decorative icon -->
 *   <OcIcon name="check" size="large" variation="success" />
 *
 *   <!-- Example of an accessible icon -->
 *   <OcIcon
 *     name="alert"
 *     fillType="line"
 *     accessibleLabel="Warning: Action required"
 *     size="medium"
 *     variation="warning"
 *   />
 *
 *   <!-- Example of a custom-colored icon -->
 *   <OcIcon
 *     name="star"
 *     color="#FFD700"
 *     size="small"
 *     variation="brand"
 *   />
 * </template>
 */

interface Props {
  name?: string
  fillType?: 'fill' | 'line' | 'none' | string
  accessibleLabel?: string
  type?: string
  size?: AvailableSizeType | string
  variation?:
    'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'brand' | 'inherit' | string
  color?: string
}
interface Emits {
  (e: 'loaded'): void
}

defineOptions({
  name: 'OcIcon',
  status: 'ready',
  release: '1.0.0'
})
const {
  name = 'info',
  fillType = 'fill',
  accessibleLabel = '',
  type = 'span',
  size = 'medium',
  variation = 'passive',
  color = ''
} = defineProps<Props>()
defineEmits<Emits>()

function sizeClass(c: string) {
  return prefix(getSizeClass(c))
}
function variationClass(c: string) {
  return prefix(c)
}
function prefix(string: string) {
  if (string !== null) {
    return `oc-icon-${string}`
  }
}
function transformSvgElement(svg: SVGElement) {
  if (accessibleLabel !== '') {
    const title = document.createElement('title')
    title.setAttribute('id', unref(svgTitleId))
    title.appendChild(document.createTextNode(accessibleLabel))
    svg.insertBefore(title, svg.firstChild)
  }
  return svg
}

const svgTitleId = computed(() => {
  return uniqueId('oc-icon-title-')
})

const nameWithFillType = computed(() => {
  const path = 'icons/'
  if (fillType.toLowerCase() === 'none') {
    return `${path}${name}.svg`
  }
  return `${path}${name}-${fillType.toLowerCase()}.svg`
})
</script>

<style lang="scss">
@mixin oc-icon-size($factor) {
  height: $oc-size-icon-default * $factor;
  max-height: $oc-size-icon-default * $factor;
  max-width: $oc-size-icon-default * $factor;
  width: $oc-size-icon-default * $factor;
}

.oc-icon {
  // SVG wrapper
  display: inline-block;
  vertical-align: baseline;

  svg {
    display: block;
  }

  &,
  > svg {
    @include oc-icon-size(1);
  }

  &-xs {
    &,
    > svg {
      @include oc-icon-size(0.5);
    }
  }

  &-s {
    &,
    > svg {
      @include oc-icon-size(0.7);
    }
  }

  &-m {
    &,
    > svg {
      @include oc-icon-size(1);
    }
  }

  &-l {
    &,
    > svg {
      @include oc-icon-size(1.5);
    }
  }

  &-xl {
    &,
    > svg {
      @include oc-icon-size(2);
    }
  }

  &-xxl {
    &,
    > svg {
      @include oc-icon-size(4);
    }
  }

  &-xxxl {
    &,
    > svg {
      @include oc-icon-size(8);
    }
  }

  &-primary > svg {
    fill: var(--oc-color-swatch-primary-default);
  }

  &-passive > svg {
    fill: var(--oc-color-swatch-passive-default);
  }

  &-warning > svg {
    fill: var(--oc-color-swatch-warning-default);
  }

  &-success > svg {
    fill: var(--oc-color-swatch-success-default);
  }

  &-danger > svg {
    fill: var(--oc-color-swatch-danger-default);
  }

  &-brand > svg {
    fill: var(--oc-color-swatch-brand-default);
  }
}
</style>
