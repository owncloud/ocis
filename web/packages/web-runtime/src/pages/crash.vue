<template>
  <div id="page-crash" class="container">
    <img class="logo" :src="logoImg" :alt="productName" />
    <div class="card">
      <h1 class="title">
        {{
          $pgettext(
            'Title of the crash page displayed when the application crashes',
            'Application crashed'
          )
        }}
      </h1>
      <p v-if="errorCode !== ''" class="error-code">
        {{ errorCode }}
      </p>
      <p>{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { useThemeStore } from '@ownclouders/web-pkg'
import { useHead } from '../composables/head'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { useGettext } from 'vue3-gettext'
import { CRASH_CODES } from '@ownclouders/web-pkg/src/errors/codes'

const route = useRoute()
const { $pgettext } = useGettext()

const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const errorCode = computed(() => route.query.code || '')

const errorMessage = computed(() => {
  switch (unref(errorCode)) {
    case CRASH_CODES.RUNTIME_BOOTSTRAP_SPACES_LOAD:
      return $pgettext(
        'An error message displayed on crash page when loading spaces during runtime bootstrap fails due to whatever reason.',
        'An error occurred while loading your spaces. Please try to reload the page. If the problem persists, seek help from your Administrator.'
      )
    default:
      return $pgettext(
        'Generic error message displayed on crash page when either no error code is provided or the error code is not known.',
        'An unknown error occurred. Please try to reload the page. If the problem persists, seek help from your Administrator.'
      )
  }
})

const productName = computed(() => currentTheme.value.common.name)
const logoImg = computed(() => currentTheme.value.logo.login)

useHead()
</script>

<style lang="scss" scoped>
.container {
  align-content: center;
  align-items: center;
  display: grid;
  gap: var(--oc-space-large);
  justify-items: center;
  min-height: 100dvh;
  text-align: center;
}

.logo {
  max-height: 12.5rem;
  max-width: 12.5rem;
}

.card {
  box-sizing: border-box;
  background-color: var(--oc-color-background-default);
  border-radius: 0.3125rem;
  color: var(--oc-color-text-default);
  padding: var(--oc-space-medium);
  margin-inline: auto;
  max-width: min(36rem, 100%);
  text-align: center;
  width: 100%;
}

.title {
  font-size: var(--oc-font-size-large);
  font-weight: var(--oc-font-weight-semibold);
  margin: 0;
}

.error-code {
  margin-top: var(--oc-space-small);
  font-size: var(--oc-font-size-xsmall);
  color: var(--oc-color-text-muted);
}
</style>
