import Notifications from '../../../../src/components/Topbar/Notifications.vue'
import { Notification } from '../../../../src/helpers/notifications'
import { mock } from 'vitest-mock-extended'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { SpaceResource } from '@ownclouders/web-client'
import { RouterLink, RouteLocationNamedRaw, RouteLocationNormalizedLoaded } from 'vue-router'
import { AxiosResponse } from 'axios'
import Avatar from '../../../../src/components/Avatar.vue'

const selectors = {
  notificationBellStub: 'notification-bell-stub',
  avatarImageStub: 'avatar-image-stub',
  noNewNotifications: '.oc-notifications-no-new',
  markAll: '.oc-notifications-mark-all',
  notificationsLoading: '.oc-notifications-loading',
  notificationItem: '.oc-notifications-item',
  notificationSubject: '.oc-notifications-subject',
  notificationMessage: '.oc-notifications-message',
  notificationLink: '.oc-notifications-link'
}
const object_id = 'normal-provider$resource-1'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  queryItemAsString: vi.fn(),
  useAppDefaults: vi.fn(),
  useVault: vi.fn(() => ({ isInVault: false }))
}))

const { useVault } = vi.mocked(await import('@ownclouders/web-pkg'))

describe('Notification component', () => {
  it('renders the notification bell and no notifications if there are none', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.find(selectors.notificationBellStub).exists()).toBeTruthy()
    expect(wrapper.find(selectors.noNewNotifications).exists()).toBeTruthy()
    expect(wrapper.find(selectors.markAll).exists()).toBeFalsy()
  })
  it('renders a set of notifications', async () => {
    const notifications = [
      mock<Notification>({
        messageRich: undefined,
        computedMessage: undefined,
        computedLink: undefined,
        object_id
      })
    ]
    const { wrapper } = getWrapper({ notifications })
    await wrapper.vm.fetchNotificationsTask.last
    expect(wrapper.find(selectors.noNewNotifications).exists()).toBeFalsy()
    expect(wrapper.findAll(selectors.notificationItem).length).toBe(notifications.length)
  })
  it('renders the loading state', async () => {
    const notifications = [mock<Notification>({ messageRich: undefined, object_id })]
    const { wrapper } = getWrapper({ notifications })
    await wrapper.vm.$nextTick()
    expect(wrapper.find(selectors.notificationsLoading).exists()).toBeTruthy()
  })
  it('marks all notifications as read', async () => {
    const notifications = [mock<Notification>({ messageRich: undefined, object_id })]
    const { wrapper, mocks } = getWrapper({ notifications })
    await wrapper.vm.fetchNotificationsTask.last
    await wrapper.find(selectors.markAll).trigger('click')
    await wrapper.vm.$nextTick()
    expect(wrapper.find(selectors.notificationItem).exists()).toBeFalsy()
    expect(mocks.$clientService.httpAuthenticated.delete).toHaveBeenCalledTimes(1)
  })
  describe('avatar', () => {
    it('loads based on the username', async () => {
      const notification = mock<Notification>({
        messageRich: undefined,
        user: 'einstein',
        object_id
      })
      const { wrapper } = getWrapper({ notifications: [notification] })
      await wrapper.vm.fetchNotificationsTask.last
      const avatarImageStub = wrapper.findComponent<typeof Avatar>(selectors.avatarImageStub)
      expect(avatarImageStub.attributes('userid')).toEqual(notification.user)
      expect(avatarImageStub.attributes('user-name')).toEqual(notification.user)
    })
    it('loads based on the rich parameters', async () => {
      const displayname = 'Albert Einstein'
      const name = 'einstein'
      const notification = mock<Notification>({
        messageRich: undefined,
        messageRichParameters: { user: { displayname, name } },
        object_id
      })
      const { wrapper } = getWrapper({ notifications: [notification] })
      await wrapper.vm.fetchNotificationsTask.last
      const avatarImageStub = wrapper.findComponent<typeof Avatar>(selectors.avatarImageStub)
      expect(avatarImageStub.attributes('userid')).toEqual(name)
      expect(avatarImageStub.attributes('user-name')).toEqual(displayname)
    })
  })
  describe('subject', () => {
    it('displays if no message given', async () => {
      const notification = mock<Notification>({
        messageRich: undefined,
        message: undefined,
        computedMessage: undefined,
        computedLink: undefined,
        object_id
      })
      const { wrapper } = getWrapper({ notifications: [notification] })
      await wrapper.vm.fetchNotificationsTask.last
      expect(wrapper.find(selectors.notificationSubject).exists()).toBeTruthy()
    })
  })
  describe('message', () => {
    it('displays simple message if messageRich not given', async () => {
      const notification = mock<Notification>({
        messageRich: undefined,
        message: 'some message',
        computedMessage: undefined,
        computedLink: undefined,
        object_id
      })
      const { wrapper } = getWrapper({ notifications: [notification] })
      await wrapper.vm.fetchNotificationsTask.last
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.notificationMessage).text()).toEqual(notification.message)
    })
    it('displays rich message and interpolates the text', async () => {
      const notification = mock<Notification>({
        messageRich: '{user} shared {resource} with you',
        messageRichParameters: {
          user: { displayname: 'Albert Einstein' },
          resource: { name: 'someFile.txt' }
        },
        computedMessage: undefined,
        computedLink: undefined,
        object_id
      })
      const { wrapper } = getWrapper({ notifications: [notification] })
      await wrapper.vm.fetchNotificationsTask.last
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.notificationMessage).text()).toEqual(
        'Albert Einstein shared someFile.txt with you'
      )
    })
    describe('escape server response', () => {
      it.each([
        ['script tags', '<script>alert("lorem")</script>hello', 'script'],
        ['img onerror (OCM share name vector)', '<img src=x onerror="alert(1)">', 'img']
      ])('strips %s from plain message', async (_label, message, selector) => {
        const notification = mock<Notification>({
          messageRich: undefined,
          message,
          computedMessage: undefined,
          computedLink: undefined,
          object_id
        })
        const { wrapper } = getWrapper({ notifications: [notification] })
        await wrapper.vm.fetchNotificationsTask.last
        await wrapper.vm.$nextTick()
        expect(wrapper.find(selector).exists()).toBe(false)
      })
      it('strips XSS from messageRich template itself', async () => {
        const notification = mock<Notification>({
          messageRich: '<img src=x onerror="alert(1)"> {user} shared with you',
          messageRichParameters: { user: { displayname: 'lorem' } },
          computedMessage: undefined,
          computedLink: undefined,
          object_id
        })
        const { wrapper } = getWrapper({ notifications: [notification] })
        await wrapper.vm.fetchNotificationsTask.last
        await wrapper.vm.$nextTick()
        expect(wrapper.find('img').exists()).toBe(false)
        expect(wrapper.find(selectors.notificationMessage).text()).toContain('lorem')
      })
      it('escapes messageRichParameters display name', async () => {
        const notification = mock<Notification>({
          messageRich: '{user} shared {resource} with you',
          messageRichParameters: {
            user: { displayname: '<img src=x onerror="alert(1)">' },
            resource: { name: 'file.txt' }
          },
          computedMessage: undefined,
          computedLink: undefined,
          object_id
        })
        const { wrapper } = getWrapper({ notifications: [notification] })
        await wrapper.vm.fetchNotificationsTask.last
        await wrapper.vm.$nextTick()
        expect(wrapper.find('img').exists()).toBe(false)
        expect(wrapper.find(selectors.notificationMessage).text()).toContain('file.txt')
      })
    })
  })
  describe('link', () => {
    it('displays if given directly', async () => {
      const notification = mock<Notification>({
        messageRich: undefined,
        computedMessage: undefined,
        computedLink: undefined,
        link: 'https://some-link.com',
        object_id
      })
      const { wrapper } = getWrapper({ notifications: [notification] })
      await wrapper.vm.fetchNotificationsTask.last
      expect(wrapper.find(selectors.notificationLink).exists()).toBeTruthy()
    })
    describe('if given via messageRichParameters', () => {
      it('renders notification as link for shares', async () => {
        const notification = mock<Notification>({
          messageRich: '{user} shared {resource} with you',
          object_type: 'share',
          messageRichParameters: {
            user: { displayname: 'Albert Einstein' },
            resource: { name: 'someFile.txt', id: '1' },
            share: { id: '1' }
          },
          computedMessage: undefined,
          computedLink: undefined,
          object_id
        })
        const { wrapper } = getWrapper({ notifications: [notification] })
        await wrapper.vm.fetchNotificationsTask.last
        await wrapper.vm.$nextTick()

        const routerLink = wrapper.findComponent<typeof RouterLink>(
          `${selectors.notificationItem} router-link-stub`
        )
        expect((routerLink.props('to') as RouteLocationNamedRaw).name).toEqual(
          'files-shares-with-me'
        )
        expect((routerLink.props('to') as RouteLocationNamedRaw).query).toEqual({
          scrollTo: notification.messageRichParameters.resource.id
        })
      })
      it('renders notification as link for spaces', async () => {
        const spaceMock = mock<SpaceResource>({
          fileId: '1',
          getDriveAliasAndItem: () => 'driveAlias',
          disabled: false
        })
        const notification = mock<Notification>({
          messageRich: '{user} added you to space {space}',
          object_type: 'storagespace',
          messageRichParameters: {
            user: { displayname: 'Albert Einstein' },
            space: { name: 'someFile.txt', id: `${spaceMock.fileId}!2` }
          },
          computedMessage: undefined,
          computedLink: undefined,
          object_id
        })
        const { wrapper } = getWrapper({ notifications: [notification], spaces: [spaceMock] })
        await wrapper.vm.fetchNotificationsTask.last
        await wrapper.vm.$nextTick()
        const routerLink = wrapper.findComponent<typeof RouterLink>(
          `${selectors.notificationItem} router-link-stub`
        )
        expect((routerLink.props('to') as RouteLocationNormalizedLoaded).params).toEqual({
          driveAliasAndItem: 'driveAlias'
        })
      })
    })
  })
  describe('vault filtering', () => {
    const vaultNotification = mock<Notification>({
      messageRich: undefined,
      computedMessage: undefined,
      computedLink: undefined,
      object_id: 'vault-provider$resource-20',
      datetime: '2024-01-01T00:00:00Z'
    })
    const driveNotification = mock<Notification>({
      messageRich: undefined,
      computedMessage: undefined,
      computedLink: undefined,
      object_id: 'other-provider$resource-22',
      datetime: '2024-01-01T00:00:01Z'
    })
    const vaultCapabilityState = {
      capabilities: { vault: { enabled: true, vault_storage_provider: 'vault-provider' } }
    }

    it('shows only vault notifications when in vault mode', async () => {
      vi.mocked(useVault).mockReturnValue({ isInVault: true })
      const { wrapper } = getWrapper({
        notifications: [vaultNotification, driveNotification],
        capabilityState: vaultCapabilityState
      })
      await wrapper.vm.fetchNotificationsTask.last
      expect(wrapper.findAll(selectors.notificationItem).length).toBe(1)
    })
    it('shows only non-vault notifications when in drive mode', async () => {
      vi.mocked(useVault).mockReturnValue({ isInVault: false })
      const { wrapper } = getWrapper({
        notifications: [vaultNotification, driveNotification],
        capabilityState: vaultCapabilityState
      })
      await wrapper.vm.fetchNotificationsTask.last
      expect(wrapper.findAll(selectors.notificationItem).length).toBe(1)
    })
    it('shows all notifications when vault is not enabled', async () => {
      vi.mocked(useVault).mockReturnValue({ isInVault: false })
      const { wrapper } = getWrapper({
        notifications: [driveNotification, driveNotification],
        capabilityState: { capabilities: { vault: { enabled: false } } }
      })
      await wrapper.vm.fetchNotificationsTask.last
      expect(wrapper.findAll(selectors.notificationItem).length).toBe(2)
    })
  })
})

function getWrapper({
  mocks = {},
  notifications = [],
  spaces = [],
  capabilityState = {}
}: {
  mocks?: Record<string, unknown>
  notifications?: Notification[]
  spaces?: SpaceResource[]
  capabilityState?: Record<string, unknown>
} = {}) {
  const localMocks = { ...defaultComponentMocks(), ...mocks }
  localMocks.$clientService.httpAuthenticated.get.mockResolvedValue(
    mock<AxiosResponse>({ data: { ocs: { data: notifications } }, headers: {} })
  )

  return {
    mocks: localMocks,
    wrapper: shallowMount(Notifications, {
      global: {
        renderStubDefaultSlot: true,
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              spacesState: { spaces },
              capabilityState
            }
          })
        ],
        mocks: localMocks,
        provide: localMocks,
        stubs: { 'avatar-image': true, OcButton: false }
      }
    })
  }
}
