import { useGettext } from 'vue3-gettext'
import translations from '../l10n/translations.json'
import { defineWebApplication } from '@ownclouders/web-pkg'
import { useExtensions } from './composables/useExtensions'
import { useApplicationFileExtensions } from './composables/useApplicationFileExtensions'

export default defineWebApplication({
  setup() {
    const { $gettext } = useGettext()
    const extensions = useExtensions()
    const applicationFileExtensions = useApplicationFileExtensions()

    return {
      appInfo: {
        name: $gettext('Password Protected Folders'),
        id: 'password-protected-folders',
        extensions: applicationFileExtensions
      },
      translations,
      extensions
    }
  }
})
