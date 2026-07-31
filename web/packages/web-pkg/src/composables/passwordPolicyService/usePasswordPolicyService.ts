import { useService } from '../service'
import { PasswordPolicyService } from '../../services/passwordPolicy'

export const usePasswordPolicyService = (): PasswordPolicyService => {
  return useService('$passwordPolicyService')
}
