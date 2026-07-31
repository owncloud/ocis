<template>
  <div ref="rootElement">
    <app-loading-spinner v-if="isLoading" />
    <template v-else>
      <p v-if="!activities.length" v-text="$gettext('No activities')" />
      <div v-else class="oc-ml-s">
        <oc-list class="timeline">
          <li v-for="activity in activities" :key="activity.id">
            <span v-html="getHtmlFromActivity(activity)" />
            <span
              class="oc-text-muted oc-text-small oc-mt-s"
              v-text="getTimeFromActivity(activity)"
            />
          </li>
        </oc-list>
        <p class="oc-text-muted oc-text-small" v-text="activitiesFooterText" />
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { computed, inject, onBeforeUnmount, onMounted, Ref, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import {
  AppLoadingSpinner,
  formatDateFromDateTime,
  useClientService,
  VisibilityObserver
} from '@ownclouders/web-pkg'
import { useTask } from 'vue-concurrency'
import { call, Resource } from '@ownclouders/web-client'
import { DateTime } from 'luxon'
import { Activity } from '@ownclouders/web-client/graph/generated'
import escape from 'lodash-es/escape'
import DOMPurify from 'dompurify'

const visibilityObserver = new VisibilityObserver()
const rootElement = ref<HTMLElement>()
const { $ngettext, current: currentLanguage } = useGettext()
const clientService = useClientService()
const resource = inject<Ref<Resource>>('resource')
const activities = ref<Activity[]>([])
const activitiesLimit = 200

const activitiesFooterText = computed(() => {
  return $ngettext(
    'Showing %{activitiesCount} activity',
    'Showing %{activitiesCount} activities',
    unref(activities).length,
    {
      activitiesCount: unref(activities).length.toString()
    }
  )
})

const loadActivitiesTask = useTask(function* (signal) {
  activities.value = yield* call(
    clientService.graphAuthenticated.activities.listActivities(
      `itemid:${unref(resource).fileId} AND limit:${activitiesLimit} AND sort:desc`,
      { signal }
    )
  )
}).restartable()

const isLoading = computed(() => {
  return loadActivitiesTask.isRunning || !loadActivitiesTask.last
})

const getHtmlFromActivity = (activity: Activity) => {
  let message = escape(activity.template.message)
  for (const [key, value] of Object.entries(activity.template.variables)) {
    const escapedValue = escape(value.displayName || value.name)

    message = message.replace(escape(`{${key}}`), `<strong>${escapedValue}</strong>`)
  }
  return DOMPurify.sanitize(message, { ALLOWED_TAGS: ['strong', 'em', 'b', 'i'], ALLOWED_ATTR: [] })
}

const getTimeFromActivity = (activity: Activity) => {
  const dateTime = DateTime.fromISO(activity.times.recordedTime)
  return formatDateFromDateTime(dateTime, currentLanguage)
}

const isVisible = ref(false)
watch(
  [resource, isVisible],
  () => {
    if (!unref(isVisible)) {
      return
    }
    loadActivitiesTask.perform()
  },
  {
    immediate: true,
    deep: true
  }
)

onMounted(() => {
  visibilityObserver.observe(unref(rootElement), {
    onEnter: () => {
      isVisible.value = true
    },
    onExit: () => {
      isVisible.value = false
    }
  })
})

onBeforeUnmount(() => {
  visibilityObserver.disconnect()
})
</script>

<style lang="scss">
.timeline {
  position: relative;
  list-style: none;
  padding: 0;
  margin: 0;

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 1.5px;
    background-color: var(--oc-color-border);
  }

  li {
    display: flex;
    flex-direction: column;
    position: relative;
    padding: 10px 20px 10px 30px;
    width: 100%;
    box-sizing: border-box;

    &::before {
      content: '';
      width: 10px;
      height: 10px;
      background-color: var(--oc-color-border);
      border-radius: 50%;
      position: absolute;
      left: -4px;
      top: 50%;
      transform: translateY(-50%);
      z-index: 1;
    }
  }
}
</style>
