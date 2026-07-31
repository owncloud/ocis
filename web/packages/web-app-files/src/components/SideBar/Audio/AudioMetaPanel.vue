<template>
  <div id="files-sidebar-panel-audio">
    <dl class="audio-data-list">
      <dt v-text="$gettext('Title')" />
      <dd data-testid="audio-panel-title" v-text="title" />
      <dt v-text="$gettext('Duration')" />
      <dd data-testid="audio-panel-duration" v-text="duration" />
      <dt v-text="$gettext('Authors')" />
      <dd data-testid="audio-panel-artist" v-text="artist" />
      <dt v-text="$gettext('Album')" />
      <dd data-testid="audio-panel-album" v-text="album" />
      <dt v-text="$gettext('Genre')" />
      <dd data-testid="audio-panel-genre" v-text="genre" />
      <dt v-text="$gettext('Year recorded')" />
      <dd data-testid="audio-panel-year" v-text="year" />
      <dt v-text="$gettext('Track')" />
      <dd data-testid="audio-panel-track" v-text="track" />
      <dt v-text="$gettext('Disc')" />
      <dd data-testid="audio-panel-disc" v-text="disc" />
    </dl>
  </div>
</template>

<script lang="ts" setup>
import { computed, inject, Ref, unref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { Duration } from 'luxon'

const resource = inject<Ref<Resource>>('resource')

const album = computed(() => {
  return unref(resource).audio.album || '-'
})

const artist = computed(() => {
  return unref(resource).audio.artist || '-'
})

const genre = computed(() => {
  return unref(resource).audio.genre || '-'
})

const title = computed(() => {
  return unref(resource).audio.title || '-'
})

const duration = computed(() => {
  const milliseconds = unref(resource).audio.duration
  if (!milliseconds) {
    return '-'
  }
  const d = Duration.fromMillis(milliseconds)
  if (d.hours > 0) {
    return d.toFormat('hh:mm:ss')
  }
  return d.toFormat('mm:ss')
})

const track = computed(() => {
  const audio = unref(resource).audio
  if (audio.track && audio.trackCount) {
    return `${audio.track} / ${audio.trackCount}`
  }
  return audio.track || '-'
})

const disc = computed(() => {
  const audio = unref(resource).audio
  if (audio.disc && audio.discCount) {
    return `${audio.disc} / ${audio.discCount}`
  }
  return audio.disc || '-'
})

const year = computed(() => {
  return unref(resource).audio.year || '-'
})
</script>

<style lang="scss">
.audio-data-list {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  dt,
  dd {
    margin-bottom: var(--oc-space-small);
  }
  dt {
    font-weight: bold;
    white-space: nowrap;
  }
  dd {
    margin-inline-start: var(--oc-space-medium);
  }
}
</style>
