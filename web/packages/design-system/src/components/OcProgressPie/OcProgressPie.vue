<template>
  <div class="oc-progress-pie" :data-fill="_fill">
    <div class="oc-progress-pie-container" />
    <span v-if="showLabel" class="oc-progress-pie-label oc-text-muted" v-text="_label" />
  </div>
</template>
<script lang="ts" setup>
import { computed } from 'vue'

/**
 * @component OcProgressPie
 * @description Displays progress in a pie chart format.
 *
 * @props {number} [progress=0] - Current value of the progress. Must be between 0 and `max`.
 * @props {number} [max=100] - Maximum value of the progress.
 * @props {boolean} [showLabel=false] - Determines if the progress label should be displayed.
 *
 * @computed {number} _fill - The calculated percentage of the progress.
 * @computed {string} _label - The label to display, formatted as either a percentage or a fraction.
 *
 * @example
 *  <!-- Basic usage -->
 *  <oc-progress-pie :progress="33" />
 *
 *  <!-- With label -->
 *  <oc-progress-pie :progress="33" show-label />
 *
 *  <!-- Custom max value -->
 *  <oc-progress-pie :progress="2" :max="4" />
 *
 *  <!-- Custom max value with label -->
 *  <oc-progress-pie :progress="4" :max="6" show-label />
 */

interface Props {
  progress?: number
  max?: number
  showLabel?: boolean
}

defineOptions({
  name: 'OcProgressPie',
  status: 'ready',
  release: '1.0.0'
})
const { progress = 0, max = 100, showLabel = false } = defineProps<Props>()
const _fill = computed(() => {
  return Math.round((100 / max) * progress)
})
const _label = computed(() => {
  if (max === 100) {
    return progress + '%'
  } else {
    return `${progress}/${max}`
  }
})
</script>
<style lang="scss">
$default-size: 64px;

.oc-progress-pie {
  float: left;
  height: $default-size;
  margin: 15px;
  position: relative;
  width: $default-size;

  *,
  *::before,
  *::after {
    box-sizing: border-box;
  }

  // Shadow
  &::after {
    border: calc($default-size / 10) solid var(--oc-color-swatch-passive-hover);
    border-radius: 50%;
    box-sizing: border-box;
    content: '';
    display: block;
    height: 100%;
    width: 100%;
  }

  &-container {
    clip: rect(0, $default-size, $default-size, calc($default-size / 2));
    height: 100%;
    left: 0;
    position: absolute;
    top: 0;
    width: 100%;

    &::before,
    &::after {
      border: calc($default-size / 10) solid var(--oc-color-swatch-brand-default);
      border-color: var(--oc-color-swatch-brand-default);
      border-radius: 50%;
      clip: rect(0, calc($default-size / 2), $default-size, 0);
      content: '';
      display: block;
      height: 100%;
      left: 0;
      position: absolute;
      top: 0;
      width: 100%;
    }
  }

  &-label {
    color: var(--oc-color-text-muted) !important;
    left: 50%;
    position: absolute;
    top: 50%;
    transform: translate(-50%, -50%);
  }
}

@for $i from 0 through 100 {
  .oc-progress-pie[data-fill='#{$i}'] {
    .oc-progress-pie-container::before {
      transform: rotate($i * 3.6deg);
    }

    @if $i <= 50 {
      .oc-progress-pie-container::after {
        display: none;
      }
    } @else {
      .oc-progress-pie-container {
        clip: rect(auto, auto, auto, auto);
      }

      .oc-progress-pie-container::after {
        transform: rotate(180deg);
      }
    }
  }
}
</style>
<docs>
```js
<section>
  <h3 class="oc-heading-divider">
    Pie shape progress
  </h3>
  <oc-progress-pie :progress="33" />
  <oc-progress-pie :progress="33" show-label/>
  <oc-progress-pie :progress="2" :max="4" />
  <oc-progress-pie :progress="4" :max="6" show-label />
</section>
```
</docs>
