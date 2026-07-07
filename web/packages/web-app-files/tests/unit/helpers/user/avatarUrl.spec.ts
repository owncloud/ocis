import { avatarUrl } from '../../../../src/helpers/user'
import { ImageDimension } from '@ownclouders/web-pkg'
import { ClientService } from '@ownclouders/web-pkg'
import { mockDeep } from 'vitest-mock-extended'
import { AxiosResponse } from 'axios'

const getDefaultOptions = () => ({
  clientService: mockDeep<ClientService>(),
  server: 'https://www.ocis.rules/',
  username: 'ocis',
  token: 'rules'
})

describe('avatarUrl', () => {
  it('throws an error', async () => {
    const defaultOptions = getDefaultOptions()
    defaultOptions.clientService.httpAuthenticated.head.mockResolvedValue({
      status: 200
    } as AxiosResponse)
    defaultOptions.clientService.ocs.signUrl.mockRejectedValue(new Error('error'))
    const avatarUrlPromise = avatarUrl(defaultOptions)
    await expect(avatarUrlPromise).rejects.toThrow(new Error('error'))
    expect(defaultOptions.clientService.httpAuthenticated.head).toHaveBeenCalledWith(
      buildUrl(defaultOptions)
    )
  })
  it('returns a signed url', async () => {
    const defaultOptions = getDefaultOptions()
    defaultOptions.clientService.httpAuthenticated.head.mockResolvedValue({
      status: 200
    } as AxiosResponse)
    defaultOptions.clientService.ocs.signUrl.mockImplementation((payload) => {
      return Promise.resolve(`${payload.url}?signed=true`)
    })
    const avatarUrlPromise = avatarUrl(defaultOptions)
    await expect(avatarUrlPromise).resolves.toBe(`${buildUrl(defaultOptions)}?signed=true`)
  })
  it('handles caching', async () => {
    const defaultOptions = getDefaultOptions()
    defaultOptions.clientService.httpAuthenticated.head.mockResolvedValue({
      status: 200
    } as AxiosResponse)
    defaultOptions.clientService.ocs.signUrl.mockImplementation((payload) =>
      Promise.resolve(payload.url)
    )

    const avatarUrlPromiseUncached = avatarUrl(defaultOptions, true)
    await expect(avatarUrlPromiseUncached).resolves.toBe(buildUrl(defaultOptions))
    expect(defaultOptions.clientService.httpAuthenticated.head).toHaveBeenCalledTimes(1)

    const avatarUrlPromiseCached = avatarUrl(defaultOptions, true)
    await expect(avatarUrlPromiseCached).resolves.toBe(buildUrl(defaultOptions))
    expect(defaultOptions.clientService.httpAuthenticated.head).toHaveBeenCalledTimes(1)

    const avatarUrlPromiseOtherSize = avatarUrl({ ...defaultOptions, size: 1 }, true)
    await expect(avatarUrlPromiseOtherSize).resolves.toBe(buildUrl({ ...defaultOptions, size: 1 }))
    expect(defaultOptions.clientService.httpAuthenticated.head).toHaveBeenCalledTimes(2)

    const avatarUrlPromiseSameUncached = avatarUrl(defaultOptions, false)
    await expect(avatarUrlPromiseSameUncached).resolves.toBe(buildUrl(defaultOptions))
    expect(defaultOptions.clientService.httpAuthenticated.head).toHaveBeenCalledTimes(3)
  })
})

const buildUrl = ({
  server,
  username,
  size
}: {
  server: string
  username: string
  size?: number
}) => [server, 'dav/avatars/', username, `/${size || ImageDimension.Avatar}.png`].join('')
