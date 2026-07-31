<template>
  <portal to="app.app-banner">
    <div v-if="isAppBannerAvailable" class="app-banner hide-desktop" :hidden="isVisible === false">
      <oc-button
        variation="brand"
        appearance="raw"
        class="app-banner-exit"
        aria-label="Close"
        @click="close"
      >
        <oc-icon name="close" size="small" />
      </oc-button>
      <div
        class="app-banner-icon"
        :style="{ 'background-image': `url('${currentTheme.appBanner.icon}')` }"
      ></div>
      <div class="info-container">
        <div>
          <div class="app-title">{{ currentTheme.appBanner.title }}</div>
          <div class="app-publisher">{{ currentTheme.appBanner.publisher }}</div>
          <div
            v-if="currentTheme.appBanner.additionalInformation !== ''"
            class="app-additional-info"
          >
            {{ $gettext(currentTheme.appBanner.additionalInformation) }}
          </div>
        </div>
      </div>
      <a
        :href="appUrl"
        target="_blank"
        class="app-banner-cta"
        rel="noopener"
        aria-label="{{ $gettext(currentTheme.appBanner.ctaText) }}"
        >{{ $gettext(currentTheme.appBanner.ctaText) }}</a
      >
    </div>
  </portal>
</template>

<script lang="ts" setup>
import { computed, ref, unref } from 'vue'
import { useRouter, useThemeStore } from '../composables'
import { buildUrl } from '../helpers/router'
import { useSessionStorage } from '@vueuse/core'
import { storeToRefs } from 'pinia'

interface Props {
  fileId: string
}
const props = defineProps<Props>()
const appBannerWasClosed = useSessionStorage('app_banner_closed', null)
const isVisible = ref<boolean>(unref(appBannerWasClosed) === null)

const router = useRouter()
const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const appBannerSettings = currentTheme.value.appBanner
const isAppBannerAvailable = computed(
  () => appBannerSettings && Object.keys(appBannerSettings).length != 0
)

const appUrl = computed(() => {
  return buildUrl(router, `/f/${props.fileId}`)
    .toString()
    .replace('https', currentTheme.value.appBanner?.appScheme)
})

const close = () => {
  isVisible.value = false
  useSessionStorage('app_banner_closed', 1)
}
</script>

<style scoped lang="scss">
.hide-desktop {
  @media (min-width: 768px) {
    display: none;
  }
}

.app-banner {
  overflow-x: hidden;
  width: 100%;
  height: 84px;
  background: #f3f3f3;
  font-family: Helvetica, sans, sans-serif;
  z-index: 5;
}

.info-container {
  position: absolute;
  top: 10px;
  left: 104px;
  display: flex;
  overflow-y: hidden;
  width: 60%;
  height: 64px;
  align-items: center;
  color: #000;
}

.app-banner-icon {
  position: absolute;
  top: 10px;
  left: 30px;
  width: 64px;
  height: 64px;
  border-radius: 15px;
  background-size: 64px 64px;
}

.app-banner-cta {
  position: absolute;
  top: 32px;
  right: 10px;
  z-index: 1;
  display: block;
  padding: 0 10px;
  min-width: 10%;
  border-radius: 5px;
  background: #f3f3f3;
  color: #1474fc;
  font-size: 18px;
  text-align: center;
  text-decoration: none;
}

.app-title {
  font-size: 14px;
}

.app-publisher,
.app-additional-info {
  font-size: 12px;
}

.app-banner-exit {
  position: absolute;
  top: 34px;
  left: 9px;
  margin: 0;
  width: 12px;
  height: 12px;
  border: 0;
  text-align: center;
  display: inline;
}
</style>
