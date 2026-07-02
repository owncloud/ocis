import { useService } from '../service'
import { AppProviderService } from '../../services'

export const useAppProviderService = (): AppProviderService => {
  return useService('$appProviderService')
}
