import { defineAsyncComponent } from 'vue'

// async component to avoid loading the huge toastjs package on page load
export const TextEditor = defineAsyncComponent(
  async () => (await import('./TextEditor.vue')).default
)
