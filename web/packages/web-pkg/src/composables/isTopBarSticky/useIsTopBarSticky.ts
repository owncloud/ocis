import { ref, onMounted, onBeforeUnmount } from 'vue'

/**
 * Determines if the topbar (including table headers) should be sticky or not.
 * With a vertical height less than 500px, the topbar should not be sticky because
 * it takes up too much space and overflows content.
 */
export const useIsTopBarSticky = () => {
  const isSticky = ref(true)

  const setIsSticky = () => {
    isSticky.value = window.innerHeight > 500
  }

  onMounted(() => {
    setIsSticky()
    window.addEventListener('resize', setIsSticky)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('resize', setIsSticky)
  })

  return { isSticky }
}
