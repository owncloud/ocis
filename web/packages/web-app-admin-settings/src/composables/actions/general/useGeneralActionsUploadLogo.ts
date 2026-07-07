import { supportedLogoMimeTypes } from '../../../defaults'
import { computed, VNodeRef, unref } from 'vue'
import { Action, useMessages } from '@ownclouders/web-pkg'
import { useAbility, useClientService, useRouter } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

export const useGeneralActionsUploadLogo = ({ imageInput }: { imageInput: VNodeRef }) => {
  const { showMessage, showErrorMessage } = useMessages()
  const { $gettext } = useGettext()
  const ability = useAbility()
  const clientService = useClientService()
  const router = useRouter()

  const uploadImage = async (ev: Event) => {
    const file = (ev.currentTarget as HTMLInputElement).files[0]

    if (!file) {
      return
    }

    if (!supportedLogoMimeTypes.includes(file.type)) {
      return showErrorMessage({ title: $gettext('The file type is unsupported') })
    }

    try {
      const httpClient = clientService.httpAuthenticated
      const formData = new FormData()
      formData.append('logo', file)
      await httpClient.post('/branding/logo', formData as never, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      showMessage({ title: $gettext('Logo was uploaded successfully') })
      setTimeout(() => {
        router.go(0)
      }, 1000)
    } catch (e) {
      console.error(e)
      showErrorMessage({
        title: $gettext('Failed to upload logo'),
        errors: [e]
      })
    }
  }

  const actions = computed((): Action[] => [
    {
      name: 'upload-logo',
      icon: 'image-add',
      label: () => {
        return $gettext('Upload logo')
      },
      isVisible: () => {
        return ability.can('update-all', 'Logo')
      },
      handler: () => {
        unref(imageInput).click()
      },
      class: 'oc-general-actions-upload-logo-trigger'
    }
  ])

  return {
    actions,
    uploadImage
  }
}
