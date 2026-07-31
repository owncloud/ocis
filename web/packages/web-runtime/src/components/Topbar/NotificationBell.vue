<template>
  <oc-button
    id="oc-notifications-bell"
    v-oc-tooltip="notificationsLabel"
    appearance="raw-inverse"
    variation="brand"
    :aria-label="notificationsLabel"
  >
    <oc-icon
      class="oc-cursor-pointer oc-flex oc-flex-middle"
      name="notification-3"
      fill-type="line"
    />
    <span
      v-if="notificationCount"
      :key="notificationCount"
      :class="{ shake: animate, badge: true }"
      v-text="notificationCountLabel"
    />
  </oc-button>
</template>
<script lang="ts">
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { ref } from 'vue'
import { watch } from 'vue'

export default {
  props: {
    notificationCount: {
      type: Number,
      default: 0
    }
  },
  setup(props) {
    const { $gettext } = useGettext()
    const animate = ref(false)
    const notificationsLabel = computed(() => $gettext('Notifications'))
    const notificationCountLabel = computed(() => {
      if (props.notificationCount > 99) {
        return '99+'
      }
      return `${props.notificationCount}`
    })

    watch(
      () => props.notificationCount,
      () => {
        animate.value = true
        setTimeout(() => {
          animate.value = false
        }, 600)
      }
    )

    return {
      animate,
      notificationsLabel,
      notificationCountLabel
    }
  }
}
</script>

<style lang="scss" scoped>
#oc-notifications-bell {
  position: relative;
  .badge {
    position: absolute;
    top: -6px;
    right: -9px;
    padding: var(--oc-space-xsmall);
    line-height: var(--oc-space-small);
    -webkit-border-radius: 30px;
    -moz-border-radius: 30px;
    border-radius: 30px;
    min-width: var(--oc-space-small);
    height: var(--oc-space-small);
    text-align: center;

    font-weight: 300;
    font-size: 11px;
    background: rgb(249, 54, 54);
    color: white;
    box-shadow: 0px 0px 2px 1px rgba(0, 0, 0, 0.5);
  }
}
.shake {
  animation: shake 0.6s cubic-bezier(0.46, 0.27, 0.59, 0.97) both;
  transform: translate3d(0, 0, 0);
}
@keyframes shake {
  10%,
  20%,
  50%,
  60% {
    transform: translate3d(-1px, 0, 0);
  }
  30%,
  40%,
  70%,
  80% {
    transform: translate3d(1px, 0, 0);
  }
}
</style>
