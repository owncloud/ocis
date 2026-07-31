import { LoadingService } from '../../services'
import { useService } from '../service'

export const useLoadingService = (): LoadingService => {
  return useService('$loadingService')
}
