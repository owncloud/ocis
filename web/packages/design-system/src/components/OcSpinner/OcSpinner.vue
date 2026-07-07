<template>
  <span :class="ocSpinnerClass" :aria-label="ariaLabel" tabindex="-1" role="img" />
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { getSizeClass, AvailableSizeType } from '../../helpers'
/**
 * OcSpinner Component
 *
 * A spinner component used to indicate loading or processing states in the application.
 *
 * @component
 * @name OcSpinner
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {string} [ariaLabel=''] - Descriptive text for screen readers. Strongly recommended unless the component is used inside a DOM node that is already `aria-hidden="true"`.
 * @prop {AvailableSizeType} [size='medium'] - The size of the spinner. Available options: `xsmall`, `small`, `medium`, `large`, `xlarge`, `xxlarge`, `xxxlarge`.
 *
 * @Accessibility
 * Making spinners accessible for screen readers
 *
 * 1. Making them (only!) programmatically focusable with `tabindex="-1"`. **This is already included in the component.**
 * 2. Giving them the [ARIA role](https://developers.google.com/web/fundamentals/accessibility/semantics-aria) of "img" (image) to change the semantics of the spinner element from something meaningless (div) to something meaningful (image). **This is already included in the component.**
 * 3. Giving them an accessible name ([explainer for the term](https://developer.paciellogroup.com/blog/2017/04/what-is-an-accessible-name/)) via aria-label/the ariaLabel prop. An element that has an accessible name, but no role is ignored by assistive technology
 * 4. During an action of the application that shows the spinner to visual users, the spinner should programmatically receive focus (see 1.), so that the accessible name of it (default: Loading, see 2.) is read out by screen reader software. After completion of the loading process the focus should be sent to a reasonable place of the newly loaded content.
 *
 * @example
 *  <oc-spinner
 *   aria-label="Loading with extra small spinner"
 *   size="xsmall" class="oc-mr-s"
 *  />
 */
interface Props {
  ariaLabel?: string
  size?: AvailableSizeType
}
defineOptions({
  name: 'OcSpinner',
  status: 'ready',
  release: '1.0.0'
})

const { ariaLabel = '', size = 'medium' } = defineProps<Props>()

const ocSpinnerClass = computed(() => {
  return ['oc-spinner', `oc-spinner-${getSizeClass(size)}`]
})
</script>

<style lang="scss">
@mixin oc-spinner-size($factor) {
  height: math.round(calc($oc-size-icon-default * $factor / 2)) * 2;
  width: math.round(calc($oc-size-icon-default * $factor / 2)) * 2;
}

.oc-spinner {
  @include oc-spinner-size(1);

  color: var(--oc-color-text-default);
  display: inline-block;

  &::after {
    animation: ball-clip-rotate 1s linear infinite;
    background: transparent;
    border: 1px solid currentColor;
    border-bottom-color: transparent;
    border-radius: 100%;
    content: '';
    display: block;
    height: 100%;
    width: 100%;
  }

  &,
  &::after {
    box-sizing: border-box;
    position: relative;
  }

  &-xs {
    @include oc-spinner-size(0.5);
  }

  &-s {
    @include oc-spinner-size(0.7);
  }

  &-search {
    @include oc-spinner-size(0.8);
  }

  &-m {
    @include oc-spinner-size(1);
  }

  &-l {
    @include oc-spinner-size(1.5);
  }

  &-xl {
    @include oc-spinner-size(2);
  }

  &-xxl {
    @include oc-spinner-size(4);
  }

  &-xxxl {
    @include oc-spinner-size(8);
  }
}

@keyframes ball-clip-rotate {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}
</style>

<docs>
```js
<div class="oc-flex oc-flex-middle">
  <oc-spinner aria-label="Loading with extra small spinner" size="xsmall" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with small spinner" size="small" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with default sized spinner" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with medium spinner" size="medium" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with large spinner" size="large" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with extra large spinner" size="xlarge" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with extra extra large spinner" size="xxlarge" class="oc-mr-s" />
  <oc-spinner aria-label="Loading with extra extra extra large spinner" size="xxxlarge" />
</div>
```
</docs>
