<template>
  <thead class="oc-thead">
    <slot />
  </thead>
</template>

<script lang="ts" setup>
defineOptions({
  name: 'OcThead',
  status: 'ready',
  release: '2.1.0'
})
</script>

<style lang="scss">
.oc-thead {
  border-bottom: 1px solid var(--oc-color-border);
}
</style>
