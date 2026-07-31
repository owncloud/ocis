import { ClientService } from '../../services'
import { useService } from '../service'

export const useClientService = (): ClientService => {
  return useService('$clientService')
}
