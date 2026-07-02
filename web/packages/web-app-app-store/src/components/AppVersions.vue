<template>
  <oc-table class="oc-width-1-1" :data="data" :fields="fields" padding-x="remove">
    <template #version="{ item }">
      v{{ item.version }}
      <oc-tag v-if="item.version === app.mostRecentVersion.version" size="small" class="oc-ml-s">
        {{ $gettext('most recent') }}
      </oc-tag>
    </template>
    <template #actions="{ item }">
      <app-actions :app="app" :version="item" />
    </template>
  </oc-table>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { App } from '../types'
import { useGettext } from 'vue3-gettext'
import AppActions from './AppActions.vue'
import { isEmpty } from 'lodash-es'

interface Props {
  app?: App
}
const { app = undefined } = defineProps<Props>()
const { $gettext } = useGettext()

const data = computed(() => {
  return (app.versions || [])
    .filter((version) => {
      if (isEmpty(version.version) || isEmpty(version.url)) {
        return false
      }
      try {
        new URL(version.url)
      } catch {
        return false
      }
      return true
    })
    .map((version) => {
      return {
        ...version,
        minOCIS: version.minOCIS ? `v${version.minOCIS}` : '-',
        id: version.version
      }
    })
})
const fields = computed(() => {
  return [
    {
      name: 'version',
      type: 'slot',
      width: 'expand',
      wrap: 'truncate',
      title: $gettext('App Version')
    },
    {
      name: 'minOCIS',
      type: 'raw',
      width: 'shrink',
      wrap: 'nowrap',
      title: $gettext('Infinite Scale Version')
    },
    {
      name: 'actions',
      type: 'slot',
      alignH: 'right',
      width: 'shrink',
      wrap: 'nowrap',
      title: ''
    }
  ]
})
</script>
