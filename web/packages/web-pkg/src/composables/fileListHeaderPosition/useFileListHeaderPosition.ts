import { ref, nextTick, onMounted, Ref } from 'vue'

export const useFileListHeaderPosition = (selector = ''): { y: Ref; refresh: () => void } => {
  const y = ref(0)
  const refresh = async (): Promise<void> => {
    await nextTick()
    const appBar = document.querySelector(selector || '#files-app-bar')
    const height = appBar ? appBar.getBoundingClientRect().height : 0

    if (y.value === height) {
      return
    }

    y.value = height
  }

  window.onresize = refresh
  onMounted(refresh)

  return { y, refresh }
}
