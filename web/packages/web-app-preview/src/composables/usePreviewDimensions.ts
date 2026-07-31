import { computed } from 'vue'

export const usePreviewDimensions = () => {
  const widths = [1024, 1280, 1920, 2160]
  const fallback = 3840
  const dimensions = computed<[number, number]>(() => {
    const width = widths.find((width) => window.innerWidth <= width) || fallback
    return [width, width]
  })

  return {
    dimensions
  }
}
