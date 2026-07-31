import { useService } from '../service'
import { ArchiverService } from '../../services'

export const useArchiverService = (): ArchiverService => {
  return useService('$archiverService')
}
