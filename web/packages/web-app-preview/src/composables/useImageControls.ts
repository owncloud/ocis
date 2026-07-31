import { PanzoomEventDetail } from '@panzoom/panzoom'
import { ref } from 'vue'

export const useImageControls = () => {
  const currentImageZoom = ref(1)
  const currentImageRotation = ref(0)
  const currentImagePositionX = ref(0)
  const currentImagePositionY = ref(0)

  const onPanZoomChanged = ({ detail }: { detail: PanzoomEventDetail }) => {
    currentImagePositionX.value = detail.x
    currentImagePositionY.value = detail.y
  }

  const resetImage = () => {
    currentImageZoom.value = 1
    currentImageRotation.value = 0
    currentImagePositionX.value = 0
    currentImagePositionY.value = 0
  }

  return {
    currentImageZoom,
    currentImageRotation,
    currentImagePositionX,
    currentImagePositionY,
    onPanZoomChanged,
    resetImage
  }
}
