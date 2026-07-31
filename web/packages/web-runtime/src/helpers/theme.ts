import { isEqual } from 'lodash-es'
import defaultTheme from '../../themes/owncloud/theme.json'
import { v4 as uuidV4 } from 'uuid'
import { ThemingConfig } from '@ownclouders/web-pkg'

export const loadTheme = async (location = '') => {
  const defaultOwnCloudTheme = {
    defaults: {
      ...defaultTheme.clients.web.defaults,
      common: defaultTheme.common
    },
    themes: defaultTheme.clients.web.themes
  }

  if (location.split('.').pop() !== 'json') {
    if (isEqual(process.env.NODE_ENV, 'development')) {
      console.info(`Theme '${location}' does not specify a json file, using default theme.`)
    }
    return defaultOwnCloudTheme
  }

  try {
    const response = await fetch(location, { headers: { 'X-Request-ID': uuidV4() } })
    if (!response.ok) {
      console.error(`Failed to load theme '${location}', invalid response. Using default theme.`)
      return defaultOwnCloudTheme
    }

    const theme = await response.json()

    try {
      const parsedTheme = ThemingConfig.parse(theme)

      return {
        defaults: {
          common: parsedTheme.common,
          ...parsedTheme.clients.web.defaults
        },
        themes: parsedTheme.clients.web.themes
      }
    } catch (error) {
      console.error(
        `Failed to load theme '${location}', invalid theme. Using default theme.`,
        error
      )
      return defaultOwnCloudTheme
    }
  } catch {
    console.error(`Failed to load theme '${location}', using default theme.`)
    return defaultOwnCloudTheme
  }
}
