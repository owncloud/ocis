import { Router } from 'vue-router'
import { useService } from '../service'

export const useRouter = (): Router => {
  return useService('$router')
}
