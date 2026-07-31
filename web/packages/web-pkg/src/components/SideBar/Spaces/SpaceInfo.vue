<template>
  <div class="space_info oc-p-s">
    <div class="space_info__body oc-text-overflow oc-flex oc-flex-middle">
      <div class="oc-mr-s">
        <oc-icon
          name="layout-grid"
          :size="resource.description ? 'large' : 'medium'"
          class="oc-display-block"
        />
      </div>
      <div>
        <h3 data-testid="space-info-name" class="oc-font-semibold oc-m-rm" v-text="resource.name" />
        <span data-testid="space-info-subtitle" v-text="resource.description" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { inject } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'

const resource = inject<SpaceResource>('resource')
</script>

<style lang="scss">
.space_info {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  grid-gap: 5px;
  background-color: var(--oc-color-background-default);
  padding: var(--oc-space-small) var(--oc-space-small) 0 var(--oc-space-small);

  &.sidebar-panel__space_info {
    border-bottom: 1px solid var(--oc-color-border);
  }

  &__body {
    text-align: left;
    font-size: var(--oc-font-size-small);

    h3 {
      font-size: var(--oc-font-size-medium);
    }
  }
}
</style>
