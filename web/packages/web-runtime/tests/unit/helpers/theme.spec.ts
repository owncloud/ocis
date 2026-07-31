import { loadTheme } from '../../../src/helpers/theme'
import defaultTheme from '../../../themes/owncloud/theme.json'
import merge from 'lodash-es/merge'
import { ThemingConfig, WebThemeConfig } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => {
  const actual = await importOriginal<any>()
  return {
    ...actual,
    ThemingConfig: {
      parse: vi.fn((arg) => arg),
      safeParse: (arg: unknown) => actual.ThemingConfig.safeParse(arg)
    }
  }
})

vi.spyOn(console, 'error').mockImplementation(() => undefined)

const defaultOwnCloudTheme = {
  defaults: {
    ...defaultTheme.clients.web.defaults,
    common: defaultTheme.common
  },
  themes: defaultTheme.clients.web.themes
}

describe('theme loading and error reporting', () => {
  it('the locally included theme should be valid', () => {
    const { success } = ThemingConfig.safeParse(defaultTheme)
    expect(success).toBeTruthy()
  })

  it('the default web theme should be valid', () => {
    const { success } = WebThemeConfig.safeParse(defaultOwnCloudTheme)
    expect(success).toBeTruthy()
  })

  it('should load the default theme if location is empty', async () => {
    const theme = await loadTheme()
    expect(theme).toMatchObject(defaultOwnCloudTheme)
  })

  it('should load the default theme if location is not a json file extension', async () => {
    const theme = await loadTheme('some_location_without_json_file_ending.xml')
    expect(theme).toMatchObject(defaultOwnCloudTheme)
  })

  it('should load the default theme if location is not found', async () => {
    vi.spyOn(global, 'fetch').mockResolvedValue(mock<Response>({ status: 404 }))
    const theme = await loadTheme('http://www.owncloud.com/unknown.json')
    expect(theme).toMatchObject(defaultOwnCloudTheme)
  })

  it('should load the default theme if location is not a valid json file', async () => {
    const customTheme = merge({}, defaultTheme, { default: { logo: { login: 'custom.svg' } } })
    vi.spyOn(global, 'fetch').mockResolvedValue(
      mock<Response>({ status: 404, json: () => Promise.resolve(customTheme) })
    )
    const theme = await loadTheme('http://www.owncloud.com/invalid.json')
    expect(theme).toMatchObject(defaultOwnCloudTheme)
  })

  it('should load the default theme if server errors', async () => {
    vi.spyOn(global, 'fetch').mockRejectedValue(new Error())
    const theme = await loadTheme('http://www.owncloud.com')
    expect(theme).toMatchObject(defaultOwnCloudTheme)
  })

  it('should load the custom theme if a custom location is given', async () => {
    const customTheme = merge({}, defaultOwnCloudTheme, {
      defaults: { logo: { login: 'custom.svg' } }
    })

    vi.spyOn(global, 'fetch').mockResolvedValue(
      mock<Response>({
        status: 404,
        json: () =>
          Promise.resolve({
            common: defaultTheme.common,
            clients: {
              web: {
                defaults: customTheme.defaults,
                themes: customTheme.themes
              }
            }
          })
      })
    )

    const theme1 = await loadTheme('http://www.owncloud.com/custom.json')
    const theme2 = await loadTheme('/custom.json')

    expect(theme1).toMatchObject(customTheme)
    expect(theme2).toMatchObject(customTheme)
  })
})
