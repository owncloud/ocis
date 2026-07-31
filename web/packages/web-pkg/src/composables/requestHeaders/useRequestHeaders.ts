import { computed } from 'vue'
import { v4 as uuidV4 } from 'uuid'
import { Auth } from '../../services'
import { useAuthStore } from '../piniaStores'
import { useGettext } from 'vue3-gettext'
import { useClientService } from '../clientService'

export const useRequestHeaders = () => {
  const authStore = useAuthStore()
  const clientService = useClientService()
  const language = useGettext()

  const headers = computed<Record<string, string>>(() => {
    const auth = new Auth({
      accessToken: authStore.accessToken,
      publicLinkToken: authStore.publicLinkToken,
      publicLinkPassword: authStore.publicLinkPassword
    })

    return {
      'Accept-Language': language.current,
      'Initiator-ID': clientService.initiatorId,
      'X-Request-ID': uuidV4(),
      ...auth.getHeaders()
    }
  })

  return { headers }
}
