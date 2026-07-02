import { onBeforeUnmount, onMounted, ref, unref } from 'vue'

export const useFullScreenMode = () => {
  const isFullScreenModeActivated = ref(false)

  const toggleFullScreenMode = () => {
    const activateFullscreen = !unref(isFullScreenModeActivated)
    isFullScreenModeActivated.value = activateFullscreen
    if (activateFullscreen) {
      if (document.documentElement.requestFullscreen) {
        document.documentElement.requestFullscreen()
      }
    } else {
      if (document.exitFullscreen) {
        document.exitFullscreen()
      }
    }
  }

  const handleFullScreenChangeEvent = () => {
    if (document.fullscreenElement === null) {
      isFullScreenModeActivated.value = false
    }
  }
  onMounted(() => {
    document.addEventListener('fullscreenchange', handleFullScreenChangeEvent)
  })
  onBeforeUnmount(() => {
    document.removeEventListener('fullscreenchange', handleFullScreenChangeEvent)
  })

  return {
    isFullScreenModeActivated,
    toggleFullScreenMode
  }
}
