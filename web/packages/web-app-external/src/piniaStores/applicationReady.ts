import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useApplicationReadyStore = defineStore('applicationReady', () => {
  const isReady = ref<boolean>(false)

  const setReady = () => {
    isReady.value = true
  }

  return {
    isReady,
    setReady
  }
})
