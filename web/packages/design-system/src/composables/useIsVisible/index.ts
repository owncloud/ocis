import { Ref, onBeforeUnmount, ref, unref, watch } from 'vue'

export const useIsVisible = ({
  target,
  mode = 'show',
  rootMargin = '100px',
  onVisibleCallback
}: {
  target: Ref<Element>
  mode?: string
  rootMargin?: string
  onVisibleCallback?: () => void
}) => {
  const isSupported = window && 'IntersectionObserver' in window
  if (!isSupported) {
    return {
      isVisible: ref(true)
    }
  }

  const isVisible = ref(false)
  const observer = new IntersectionObserver(
    (intersectionObserverEntries: IntersectionObserverEntry[]) => {
      /**
       * In some edge cases intersectionObserverEntries contains 2 entries with the first one having wrong rootBounds.
       * This happens for some reason when the table is being re-sorted immediately after being rendered.
       * Therefore we always check the last entry for isIntersecting.
       */
      const isIntersecting = intersectionObserverEntries.at(-1).isIntersecting

      isVisible.value = isIntersecting
      if (unref(isVisible) && onVisibleCallback) {
        onVisibleCallback()
      }

      /**
       * if given mode is `showHide` we need to keep the observation alive.
       */
      if (mode === 'showHide') {
        return
      }
      /**
       * if the mode is `show` which is the default, the implementation needs to unsubscribe the target from the observer
       */
      if (!isVisible.value) {
        return
      }

      observer.unobserve(target.value)
    },
    {
      rootMargin
    }
  )

  watch(target, () => {
    observer.observe(target.value)
  })

  onBeforeUnmount(() => observer.disconnect())

  return {
    isVisible
  }
}
