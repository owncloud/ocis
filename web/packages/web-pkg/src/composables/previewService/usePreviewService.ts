import { useService } from '../service'
import { PreviewService } from '../../services/preview'

export const usePreviewService = (): PreviewService => {
  return useService('$previewService')
}
