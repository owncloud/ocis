<template>
  <li class="oc-sidebar-nav-item oc-pb-xs oc-px-s" :aria-current="active ? 'page' : null">
    <oc-button
      v-oc-tooltip="toolTip"
      :type="handler ? 'button' : 'router-link'"
      :appearance="active ? 'raw-inverse' : 'raw'"
      :variation="active ? 'primary' : 'passive'"
      :class="['oc-sidebar-nav-item-link', 'oc-oc-width-1-1', { active: active }]"
      :data-nav-id="index"
      :data-nav-name="navName"
      v-bind="attrs"
    >
      <span class="oc-flex">
        <oc-icon :name="icon" :fill-type="fillType" variation="inherit" />
        <span class="oc-ml-m text" :class="{ 'text-invisible': collapsed }" v-text="name" />
      </span>
    </oc-button>
  </li>
</template>
<script lang="ts">
import { computed, defineComponent, PropType } from 'vue'
import { RouteLocationRaw } from 'vue-router'

export default defineComponent({
  props: {
    name: {
      type: String,
      required: true
    },
    index: {
      type: String,
      required: true
    },
    active: {
      type: Boolean,
      required: false,
      default: false
    },
    target: {
      type: [String, Object] as PropType<RouteLocationRaw>,
      required: false,
      default: null
    },
    icon: {
      type: String,
      required: true
    },
    fillType: {
      type: String,
      required: false,
      default: 'fill'
    },
    collapsed: {
      type: Boolean,
      required: false,
      default: false
    },
    handler: {
      type: Function as PropType<() => void>,
      required: false,
      default: null
    }
  },
  setup(props) {
    const attrs = computed(() => {
      return {
        ...(props.handler && { onClick: props.handler }),
        ...(props.target && { to: props.target })
      }
    })

    return { attrs }
  },
  computed: {
    navName() {
      if (this.target) {
        return this.$router?.resolve(this.target, this.$route)?.name || 'route.name'
      }
      return this.name
    },
    toolTip() {
      const value = this.collapsed
        ? this.$gettext('Navigate to %{ pageName } page', {
            pageName: this.name
          })
        : ''
      return {
        content: value,
        placement: 'right',
        arrow: false
      }
    }
  }
})
</script>

<style lang="scss">
.oc-sidebar-nav-item-link {
  position: relative;
  align-items: center !important;
  display: flex !important;
  justify-content: space-between !important;
  padding: var(--oc-space-small) !important;
  border-radius: 5px;
  white-space: nowrap;
  user-select: none;

  .oc-tag {
    color: var(--oc-color-text-default);
    background-color: var(--oc-color-background-highlight);
  }
  .text {
    opacity: 1;
  }
  .text-invisible {
    opacity: 0 !important;
    transition: 0s;
  }

  &:hover:not(.active) {
    background-color: var(--oc-color-swatch-primary-hover) !important;
    color: var(--oc-color-swatch-primary-contrast) !important;

    .oc-icon > svg {
      fill: var(--oc-color-swatch-primary-contrast);
    }
  }

  &:hover,
  &:focus {
    text-decoration: none !important;
  }
  &.active {
    overflow: hidden;

    &:hover .oc-icon > svg {
      fill: var(--oc-color-swatch-primary-contrast);
    }
  }
}
</style>
