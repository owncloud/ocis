import { useGettext } from 'vue3-gettext'
import { computed } from 'vue'
import { useAbility, useClientService, useMessages, useRouter } from '@ownclouders/web-pkg'
import { Action } from '@ownclouders/web-pkg'

export const useGeneralActionsResetLogo = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const { $gettext } = useGettext()
  const ability = useAbility()
  const clientService = useClientService()
  const router = useRouter()

  const handler = async () => {
    try {
      const httpClient = clientService.httpAuthenticated
      await httpClient.delete('/branding/logo')
      showMessage({ title: $gettext('Logo was reset successfully. Reloading page...') })
      setTimeout(() => {
        router.go(0)
      }, 1000)
    } catch (e) {
      console.error(e)
      showErrorMessage({
        title: $gettext('Failed to reset logo'),
        errors: [e]
      })
    }
  }

  const actions = computed((): Action[] => [
    {
      name: 'reset-logo',
      icon: 'restart',
      label: () => {
        return $gettext('Reset logo')
      },
      isVisible: () => {
        return ability.can('update-all', 'Logo')
      },
      handler,
      class: 'oc-general-actions-reset-logo-trigger'
    }
  ])

  return {
    actions
  }
}
