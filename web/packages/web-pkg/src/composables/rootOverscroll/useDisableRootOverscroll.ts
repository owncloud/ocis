import { onBeforeUnmount, onMounted } from 'vue'

/**
 * Disables the vertical overscroll (rubber-band bounce, scroll chaining) of the
 * page root while the calling component is mounted.
 */
export const useDisableRootOverscroll = () => {
  let previousValue = ''

  onMounted(() => {
    const rootStyle = document.documentElement.style
    previousValue = rootStyle.getPropertyValue('overscroll-behavior-y')
    rootStyle.setProperty('overscroll-behavior-y', 'none')
  })

  onBeforeUnmount(() => {
    document.documentElement.style.setProperty('overscroll-behavior-y', previousValue)
  })
}
