<template>
  <div id="oc-notifications">
    <notification-bell :notification-count="notifications.length" />
    <oc-drop
      id="oc-notifications-drop"
      drop-id="notifications-dropdown"
      toggle="#oc-notifications-bell"
      mode="click"
      :options="{ pos: 'bottom-right', delayHide: 0 }"
      class="oc-overflow-auto"
    >
      <div class="oc-flex oc-flex-between oc-flex-middle oc-mb-s">
        <span class="oc-text-bold oc-text-large oc-m-rm" v-text="$gettext('Notifications')" />
        <oc-button
          v-if="notifications.length"
          class="oc-notifications-mark-all"
          appearance="raw"
          @click="deleteNotificationsTask.perform(notifications.map((n) => n.notification_id))"
          ><span v-text="$gettext('Mark all as read')"
        /></oc-button>
      </div>
      <hr />
      <div class="oc-position-relative">
        <div v-if="loading" class="oc-notifications-loading">
          <div class="oc-notifications-loading-background oc-width-1-1 oc-height-1-1" />
          <oc-spinner class="oc-notifications-loading-spinner" size="large" />
        </div>
        <span
          v-if="!notifications.length"
          class="oc-notifications-no-new"
          tabindex="0"
          v-text="$gettext('Nothing new')"
        />
        <oc-list v-else>
          <li v-for="(el, index) in notifications" :key="index" class="oc-notifications-item">
            <component
              :is="el.computedLink ? 'router-link' : 'div'"
              class="oc-flex oc-flex-middle oc-my-xs"
              :to="el.computedLink"
            >
              <avatar-image
                class="oc-mr-m"
                :width="32"
                :userid="el.messageRichParameters?.user?.name || el.user"
                :user-name="el.messageRichParameters?.user?.displayname || el.user"
              />
              <div>
                <div v-if="!el.message && !el.messageRich" class="oc-notifications-subject">
                  <span v-text="el.subject" />
                </div>
                <div v-if="el.computedMessage" class="oc-notifications-message">
                  <span v-bind="{ innerHTML: el.computedMessage }" />
                </div>
                <div
                  v-if="el.link && el.object_type !== 'local_share'"
                  class="oc-notifications-link"
                >
                  <a :href="el.link" target="_blank" v-text="el.link" />
                </div>
                <div v-if="el.datetime" class="oc-text-small oc-text-muted oc-mt-xs">
                  <span
                    v-oc-tooltip="formatDate(el.datetime)"
                    tabindex="0"
                    v-text="formatDateRelative(el.datetime)"
                  />
                </div>
              </div>
            </component>
            <hr v-if="index + 1 !== notifications.length" class="oc-my-m" />
          </li>
        </oc-list>
      </div>
    </oc-drop>
  </div>
</template>
<script lang="ts">
import { computed, onMounted, onUnmounted, ref, unref } from 'vue'
import DOMPurify from 'dompurify'
import isEmpty from 'lodash-es/isEmpty'
import escape from 'lodash-es/escape'
import {
  useCapabilityStore,
  useSpacesStore,
  createFileRouteOptions,
  formatDateFromISO,
  formatRelativeDateFromISO,
  useClientService,
  useVault
} from '@ownclouders/web-pkg'
import NotificationBell from './NotificationBell.vue'
import { Notification } from '../../helpers/notifications'
import { useGettext } from 'vue3-gettext'
import { useTask } from 'vue-concurrency'
import { MESSAGE_TYPE } from '@ownclouders/web-client/sse'
import { call } from '@ownclouders/web-client'
import { AxiosHeaders } from 'axios'

const POLLING_INTERVAL = 30000

export default {
  components: {
    NotificationBell
  },
  setup() {
    const spacesStore = useSpacesStore()
    const capabilityStore = useCapabilityStore()
    const clientService = useClientService()
    const language = useGettext()

    const rawNotifications = ref<Notification[]>([])
    const notificationsInterval = ref()
    const { isInVault } = useVault()

    const isVaultNotification = (notification: Notification) =>
      notification.object_id?.split('$')[0] === capabilityStore.vaultStorageProvider

    const notifications = computed(() => {
      if (!capabilityStore.vaultEnabled) {
        return unref(rawNotifications)
      }
      if (isInVault) {
        return unref(rawNotifications).filter(isVaultNotification)
      }
      return unref(rawNotifications).filter(
        (notification: Notification) => !isVaultNotification(notification)
      )
    })

    const loading = computed(() => {
      return fetchNotificationsTask.isRunning || deleteNotificationsTask.isRunning
    })

    const formatDate = (date: string) => {
      return formatDateFromISO(date, language.current)
    }
    const formatDateRelative = (date: string) => {
      return formatRelativeDateFromISO(date, language.current)
    }

    const messageParameters = [
      { name: 'user', labelAttribute: 'displayname' },
      { name: 'resource', labelAttribute: 'name' },
      { name: 'space', labelAttribute: 'name' },
      { name: 'virus', labelAttribute: 'name' }
    ]
    const sanitizeAllowList = { ALLOWED_TAGS: ['strong', 'em', 'b', 'i'], ALLOWED_ATTR: [] }
    const getMessage = ({ message, messageRich, messageRichParameters }: Notification): string => {
      if (messageRich && !isEmpty(messageRichParameters)) {
        let interpolatedMessage = escape(messageRich)
        for (const param of messageParameters) {
          const placeholder = escape(`{${param.name}}`)
          if (interpolatedMessage.includes(placeholder)) {
            const richParam = messageRichParameters[param.name] ?? undefined
            if (!richParam) {
              return DOMPurify.sanitize(escape(message), sanitizeAllowList)
            }
            const label = richParam[param.labelAttribute] ?? undefined
            if (!label) {
              return DOMPurify.sanitize(escape(message), sanitizeAllowList)
            }
            interpolatedMessage = interpolatedMessage.replace(
              placeholder,
              `<strong>${escape(label)}</strong>`
            )
          }
        }
        return DOMPurify.sanitize(interpolatedMessage, sanitizeAllowList)
      }
      return DOMPurify.sanitize(escape(message), sanitizeAllowList)
    }
    const getLink = ({ messageRichParameters, object_type }: Notification) => {
      if (!messageRichParameters) {
        return null
      }
      if (object_type === 'share') {
        return {
          name: 'files-shares-with-me',
          ...(!!messageRichParameters?.resource.id && {
            query: { scrollTo: messageRichParameters.resource.id }
          })
        }
      }
      if (object_type === 'storagespace' && messageRichParameters?.space?.id) {
        const space = spacesStore.spaces.find(
          (s) => s.fileId === messageRichParameters?.space?.id.split('!')[0] && !s.disabled
        )
        if (space) {
          return {
            name: 'files-spaces-generic',
            ...createFileRouteOptions(space, { path: '', fileId: space.fileId })
          }
        }
      }
      return null
    }

    const fetchNotificationsTask = useTask(function* (signal) {
      try {
        const response = yield* call(
          clientService.httpAuthenticated.get<{ ocs: { data: Notification[] } }>(
            'ocs/v2.php/apps/notifications/api/v1/notifications',
            { signal }
          )
        )

        if ((response.headers as AxiosHeaders).get('Content-Length') === '0') {
          return
        }

        const {
          ocs: { data = [] }
        } = response.data
        const sorted = data?.sort((a, b) => b.datetime.localeCompare(a.datetime)) || []
        sorted.forEach((notification) => setAdditionalNotificationData(notification))
        rawNotifications.value = sorted
      } catch (e) {
        console.error(e)
      }
    }).restartable()

    const deleteNotificationsTask = useTask(function* (signal, ids) {
      try {
        yield clientService.httpAuthenticated.delete(
          'ocs/v2.php/apps/notifications/api/v1/notifications',
          { data: { ids } },
          { signal }
        )
      } catch (e) {
        console.error(e)
      } finally {
        rawNotifications.value = unref(rawNotifications).filter(
          (n) => !ids.includes(n.notification_id)
        )
      }
    }).restartable()

    const setAdditionalNotificationData = (notification: Notification) => {
      notification.computedMessage = getMessage(notification)
      notification.computedLink = getLink(notification)
    }

    const onSSENotificationEvent = (event: MessageEvent) => {
      try {
        const notification = JSON.parse(event.data) as Notification
        if (!notification || !notification.notification_id) {
          return
        }
        setAdditionalNotificationData(notification)
        rawNotifications.value = [notification, ...unref(rawNotifications)]
      } catch {
        console.error('Unable to parse sse notification data')
      }
    }

    onMounted(() => {
      fetchNotificationsTask.perform()
      if (unref(capabilityStore.supportSSE)) {
        clientService.sseAuthenticated.addEventListener(
          MESSAGE_TYPE.NOTIFICATION,
          onSSENotificationEvent
        )
      } else {
        notificationsInterval.value = setInterval(() => {
          fetchNotificationsTask.perform()
        }, POLLING_INTERVAL)
      }
    })

    onUnmounted(() => {
      if (unref(capabilityStore.supportSSE)) {
        clientService.sseAuthenticated.removeEventListener(
          MESSAGE_TYPE.NOTIFICATION,
          onSSENotificationEvent
        )
      } else {
        clearInterval(unref(notificationsInterval))
      }
    })

    return {
      notifications,
      fetchNotificationsTask,
      loading,
      deleteNotificationsTask,
      formatDate,
      formatDateRelative,
      getMessage,
      getLink
    }
  }
}
</script>

<style lang="scss" scoped>
#oc-notifications-drop {
  width: 400px;
  max-width: 100%;
  max-height: 400px;
}

.oc-notifications {
  &-item {
    > a {
      color: var(--oc-color-text-default);
    }
  }

  &-loading {
    * {
      position: absolute;
    }

    &-background {
      background-color: var(--oc-color-background-secondary);
      opacity: 0.6;
    }

    &-spinner {
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      opacity: 1;
    }
  }

  &-actions {
    button:not(:last-child) {
      margin-right: var(--oc-space-small);
    }
  }

  &-link {
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
    width: 300px;
  }
}
</style>
