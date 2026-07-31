<template>
  <span :class="['oc-remote-icon', getSizeClass(size), variation]" :style="style" />
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { type AvailableSizeType, getSizeClass } from '../../helpers'

defineOptions({
  name: 'OcRemoteIcon'
})

const {
  src,
  size = 'medium',
  variation = 'passive'
} = defineProps<{
  src: string
  size?: AvailableSizeType
  variation?: 'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'brand' | 'inherit'
}>()

const style = computed(() => ({
  '--oc-remote-icon-src': `url(${src})`
}))
</script>

<style lang="scss" scoped>
@mixin oc-icon-size($factor) {
  height: $oc-size-icon-default * $factor;
  max-height: $oc-size-icon-default * $factor;
  max-width: $oc-size-icon-default * $factor;
  width: $oc-size-icon-default * $factor;
}

.oc-remote-icon {
  @include oc-icon-size(1);

  background-color: currentColor;

  -webkit-mask: var(--oc-remote-icon-src) no-repeat center / contain;
  mask: var(--oc-remote-icon-src) no-repeat center / contain;

  &.xs {
    @include oc-icon-size(0.5);
  }

  &.s {
    @include oc-icon-size(0.7);
  }

  &.m {
    @include oc-icon-size(1);
  }

  &.l {
    @include oc-icon-size(1.5);
  }

  &.xl {
    @include oc-icon-size(2);
  }

  &.xxl {
    @include oc-icon-size(4);
  }

  &.xxxl {
    @include oc-icon-size(8);
  }

  &.primary > svg {
    background-color: var(--oc-color-swatch-primary-default);
  }

  &.passive > svg {
    background-color: var(--oc-color-swatch-passive-default);
  }

  &.warning > svg {
    background-color: var(--oc-color-swatch-warning-default);
  }

  &.success > svg {
    background-color: var(--oc-color-swatch-success-default);
  }

  &.danger > svg {
    background-color: var(--oc-color-swatch-danger-default);
  }

  &.brand > svg {
    background-color: var(--oc-color-swatch-brand-default);
  }
}
</style>
