<template>
  <button class="skip-button" @click="skipToTarget"><slot /></button>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  props: {
    /*
     * The element to focus and to skip to
     */
    target: {
      type: String,
      required: true
    }
  },
  methods: {
    skipToTarget() {
      const targetElement = document.getElementById(this.target)
      if (!targetElement) {
        return
      }
      targetElement.setAttribute('tabindex', '-1')
      targetElement.focus()
      targetElement.scrollIntoView()
    }
  }
})
</script>

<style scoped>
.skip-button {
  position: absolute;
  top: -100px;
  left: 0;
  z-index: 6;
  -webkit-appearance: none;
  border: none;
  background-color: var(--oc-color-swatch-brand-default);
  color: var(--oc-color-swatch-primary-contrast);
  font: inherit;
  padding: 0.25em 0.5em;
}

.skip-button:focus {
  top: 0;
  outline: none;
  border: 1px dashed white;
}
</style>
