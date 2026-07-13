import { AppMenuItemExtension, WebThemeType } from '@ownclouders/web-pkg'
import { Language } from 'vue3-gettext'

export const buildStaticAppMenuItems = (
  common: WebThemeType['common'],
  $pgettext: Language['$pgettext']
): AppMenuItemExtension[] => {
  const items: AppMenuItemExtension[] = []

  const softwareLicenseUrl = common?.urls?.softwareLicense
  const helpPageUrl = common?.urls?.helpPage

  if (softwareLicenseUrl) {
    items.push({
      id: 'app.runtime.header.app-menu.software-license',
      type: 'appMenuItem',
      label: () =>
        $pgettext(
          'Apps menu: link label; opens the software license information page.',
          'Software License Information'
        ),
      icon: 'scales',
      priority: 900,
      url: softwareLicenseUrl
    })
  }

  if (helpPageUrl) {
    items.push({
      id: 'app.runtime.header.app-menu.help-page',
      type: 'appMenuItem',
      label: () => $pgettext('Apps menu: link label; opens the help pages.', 'Help Pages'),
      icon: 'question',
      priority: 910,
      url: helpPageUrl
    })
  }

  return items
}
