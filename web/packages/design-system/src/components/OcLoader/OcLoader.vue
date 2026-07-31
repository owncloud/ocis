<template>
  <div :class="['oc-loader', { 'oc-loader-flat': flat }]">
    <span class="oc-invisible-sr" data-testid="oc-loader-label" v-text="ariaLabel" />
  </div>
</template>

<script lang="ts" setup>
/**
 * OcLoader Component
 *
 * A loader component that provides visual feedback for ongoing actions.
 *
 * @component
 * @name OcLoader
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {string} [ariaLabel='Loading'] - Descriptive text to be read by screen readers.
 * @prop {boolean} [flat=false] - Removes border radius and shrinks the height when set to true.
 *
 * @example
 * <h3>Default style</h3>
 * <div>
 *   <oc-loader />
 * </div>
 *
 * <h3>Flat style</h3>
 * <div>
 *   <oc-loader :flat="true" />
 * </div>
 */

interface Props {
  ariaLabel?: string
  flat?: boolean
}
defineOptions({
  name: 'OcLoader',
  status: 'ready',
  release: '1.0.0'
})

const { ariaLabel = 'Loading', flat = false } = defineProps<Props>()
</script>

<style lang="scss">
.oc-loader {
  -webkit-appearance: none;
  -moz-appearance: none;
  background-color: #f8f8f8;
  border: 0;
  border-radius: 500px;
  display: block;
  height: 15px;
  margin-top: 20px;
  margin-bottom: 20px;
  overflow: hidden;
  vertical-align: baseline;
  width: 100%;
  position: relative;

  &-flat {
    border-radius: 0 !important;
    height: 5px !important;
  }

  &::after {
    background: var(--oc-color-text-muted);
    content: '';
    height: 100%;
    width: 0;
    display: block;
    position: absolute;

    animation: {
      duration: 1.4s;
      iteration-count: infinite;
      name: oc-loader;
    }
  }
}

@keyframes oc-loader {
  0% {
    left: 0;
    width: 0;
  }

  50% {
    left: 0;
    width: 66%;
  }

  100% {
    left: 100%;
    width: 10%;
  }
}
</style>

<docs>
```js
<h3 class="oc-heading-divider">
  Default style
</h3>
<div>
  <oc-loader />
</div>

<h3 class="oc-heading-divider">
  Flat style
</h3>
<div>
  <oc-loader :flat="true" />
</div>
```
</docs>
