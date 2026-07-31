<template>
  <div
    class="no-content-message oc-flex oc-flex-column oc-flex-center oc-flex-middle oc-text-center"
  >
    <oc-icon
      :name="props.icon"
      type="div"
      size="xxlarge"
      :fill-type="props.iconFillType"
      class="oc-mb-m"
    />
    <div class="oc-text-muted oc-text-xlarge">
      <slot name="message" />
    </div>
    <div class="oc-text-muted">
      <slot name="callToAction" />
    </div>
  </div>
</template>

<script lang="ts" setup>
interface Props {
  icon: string
  iconFillType?: string
}
const props = withDefaults(defineProps<Props>(), {
  iconFillType: 'fill'
})
</script>
<style>
.no-content-message {
  height: 75dvh;
}

.space-frontpage .no-content-message {
  height: 55dvh;
}
</style>
