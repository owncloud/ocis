<template>
  <div
    :class="classes"
    :aria-valuemax="max"
    :aria-valuenow="value"
    aria-busy="true"
    aria-valuemin="0"
    role="progressbar"
  >
    <div v-if="!indeterminate" class="oc-progress-current" :style="{ width: progressValue }"></div>
    <div v-else class="oc-progress-indeterminate">
      <div class="oc-progress-indeterminate-first"></div>
      <div class="oc-progress-indeterminate-second"></div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

/**
 * OcProgress Component
 *
 * A progress bar component that visually represents the completion status of a task or process.
 *
 * @component
 * @name OcProgress
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {number} [value=0] - The current progress value.
 * @prop {number} [max] - The maximum value for the progress bar. If not provided, the progress bar will not display a percentage.
 * @prop {'default' | 'small'} [size='default'] - The size of the progress bar. Can be 'default' or 'small'.
 * @prop {'primary' | 'passive' | 'success' | 'warning' | 'danger'} [variation='primary'] - The visual variation of the progress bar.
 * @prop {boolean} [indeterminate=false] - If true, the progress bar will display an indeterminate animation.
 * @example
 *  <!-- Default progress bar -->
 *  <OcProgress :value="50" :max="100" />
 *
 *  <!-- Small progress bar -->
 *  <OcProgress :value="30" :max="100" size="small" />
 *
 *  <!-- Indeterminate progress bar -->
 *  <OcProgress indeterminate />
 *
 *  <!-- Progress bar with variation -->
 *  <OcProgress :value="70" :max="100" variation="success" />
 *
 */

interface Props {
  value?: number
  max?: number
  size?: 'default' | 'small'
  variation?: 'primary' | 'passive' | 'success' | 'warning' | 'danger'
  indeterminate?: boolean
}
defineOptions({
  name: 'OcProgress',
  status: 'ready',
  release: '1.0.0'
})
const {
  value = 0,
  max = undefined,
  size = 'default',
  variation = 'primary',
  indeterminate = false
} = defineProps<Props>()

const classes = computed(() => {
  return `oc-progress oc-progress-${size} oc-progress-${variation}`
})
const progressValue = computed(() => {
  if (!max) {
    return '-'
  }
  const num = (value / max) * 100
  return `${num}%`
})
</script>

<style lang="scss">
$progress-height: 15px !default;
$progress-height-small: 5px !default;

.oc-progress {
  background-color: var(--oc-color-input-border);
  display: block;
  height: $progress-height;
  // Add the correct vertical alignment in Chrome, Firefox, and Opera.
  width: 100%;
  position: relative;
  overflow-x: hidden;

  &-small {
    height: $progress-height-small;
  }
  &-current {
    height: 100%;
    position: absolute;
    transition: width 0.5s;
  }
  &-indeterminate div {
    height: 100%;
    position: absolute;
  }
  &-indeterminate-first {
    animation-duration: 2s;
    animation-name: indeterminate-first;
    animation-iteration-count: infinite;
  }
  &-indeterminate-second {
    animation-duration: 2s;
    animation-delay: 0.5s;
    animation-name: indeterminate-second;
    animation-iteration-count: infinite;
  }

  @keyframes indeterminate-first {
    from {
      left: -10%;
      width: 10%;
    }
    to {
      left: 120%;
      width: 100%;
    }
  }

  @keyframes indeterminate-second {
    from {
      left: -100%;
      width: 80%;
    }
    to {
      left: 110%;
      width: 10%;
    }
  }

  &-primary &-current,
  &-primary &-indeterminate div {
    background-color: var(--oc-color-swatch-primary-default);
  }
  &-success &-current,
  &-success &-indeterminate div {
    background-color: var(--oc-color-swatch-success-default);
  }
  &-warning &-current,
  &-warning &-indeterminate div {
    background-color: var(--oc-color-swatch-warning-default);
  }
  &-danger &-current,
  &-danger &-indeterminate div {
    background-color: var(--oc-color-swatch-danger-default);
  }
}
</style>
