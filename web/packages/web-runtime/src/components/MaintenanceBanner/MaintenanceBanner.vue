<template>
  <div class="maintenance-banner">
    <button class="maintenance-banner-text" @click="openDetailsModal">
      {{ $gettext('System maintenance in progress. Click here to read IMPORTANT details.') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useModals } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { MaintenanceDetailsModal } from './'

const { $gettext } = useGettext()
const { dispatchModal } = useModals()

const openDetailsModal = () => {
  dispatchModal({
    title: $gettext('System maintenance in progress'),
    customComponent: MaintenanceDetailsModal,
    cancelText: $gettext('Close'),
    hideConfirmButton: true,
    variation: 'warning'
  })
}
</script>

<style scoped>
.maintenance-banner {
  align-items: center;
  background-color: var(--oc-color-swatch-warning-default);
  box-sizing: border-box;
  display: flex;
  justify-content: center;
}

.maintenance-banner-text {
  background-color: var(--oc-color-swatch-warning-default);
  border: none;
  color: var(--oc-color-swatch-warning-contrast);
  display: block;
  height: 100%;
  margin: 0;
  padding: 0.5rem 1rem;
  text-align: center;
  width: 100%;
}
</style>
