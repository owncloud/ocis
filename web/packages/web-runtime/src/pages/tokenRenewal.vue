<template>
  <router-view />
</template>
<script lang="ts">
import { defineComponent } from 'vue'

/**
 * This page gets rendered in the iFrame for the silent token renewal to take place.
 */
export default defineComponent({
  name: 'TokenRenewal'
})
</script>
