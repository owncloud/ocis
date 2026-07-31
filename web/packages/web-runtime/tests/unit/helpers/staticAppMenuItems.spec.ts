import { buildStaticAppMenuItems } from '../../../src/helpers/staticAppMenuItems'
import { WebThemeType } from '@ownclouders/web-pkg'

const $pgettext = (context: string, msgid: string) => msgid

describe('buildStaticAppMenuItems', () => {
  it('returns an empty array when neither url is set', () => {
    const items = buildStaticAppMenuItems({ urls: {} } as WebThemeType['common'], $pgettext)
    expect(items).toHaveLength(0)
  })
  it('returns one item when only softwareLicense is set', () => {
    const items = buildStaticAppMenuItems(
      { urls: { softwareLicense: 'https://example.com/license' } } as WebThemeType['common'],
      $pgettext
    )
    expect(items).toHaveLength(1)
    expect(items[0].url).toEqual('https://example.com/license')
    expect(items[0].label()).toEqual('Software License Information')
  })
  it('returns one item when only helpPage is set', () => {
    const items = buildStaticAppMenuItems(
      { urls: { helpPage: 'https://example.com/help' } } as WebThemeType['common'],
      $pgettext
    )
    expect(items).toHaveLength(1)
    expect(items[0].url).toEqual('https://example.com/help')
    expect(items[0].label()).toEqual('Help Pages')
  })
  it('returns two items when both urls are set, sorted after real apps by priority', () => {
    const items = buildStaticAppMenuItems(
      {
        urls: {
          softwareLicense: 'https://example.com/license',
          helpPage: 'https://example.com/help'
        }
      } as WebThemeType['common'],
      $pgettext
    )
    expect(items).toHaveLength(2)
    items.forEach((item) => {
      expect(item.priority).toBeGreaterThanOrEqual(900)
      expect(item.type).toEqual('appMenuItem')
    })
  })
})
