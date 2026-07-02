<template>
  <oc-button
    v-oc-tooltip="
      $pgettext(
        'Top bar: tooltip for the “Universal Access” icon; imperative sentence that opens a menu with accessibility options.',
        'Open accessibility options'
      )
    "
    appearance="raw-inverse"
    variation="brand"
    :aria-label="
      $pgettext(
        'Top bar: label for the accessibility menu trigger; Title Case.',
        'Universal Access'
      )
    "
  >
    <oc-remote-icon :src="currentTheme.icons?.universalAccess" />
  </oc-button>
  <oc-drop ref="menu">
    <oc-list class="universal-access-dropdown-list">
      <li v-if="currentTheme.common.urls.universalAccessEasyLanguage">
        <oc-button
          appearance="raw"
          type="a"
          :href="currentTheme.common?.urls?.universalAccessEasyLanguage"
          target="_blank"
          rel="noopener noreferrer"
        >
          <oc-remote-icon
            v-if="currentTheme.icons?.universalAccessEasyLanguage"
            :src="currentTheme.icons?.universalAccessEasyLanguage"
          />
          <span
            v-text="
              $pgettext(
                'Accessibility menu: entry that links to a simplified/plain-language help or documentation page; Title Case.',
                'Easy Language'
              )
            "
          />
        </oc-button>
      </li>
      <li v-if="currentTheme.common.urls.universalAccessSignLanguage">
        <oc-button
          appearance="raw"
          type="a"
          :href="currentTheme.common?.urls?.universalAccessSignLanguage"
          target="_blank"
          rel="noopener noreferrer"
        >
          <oc-remote-icon
            v-if="currentTheme.icons?.universalAccessSignLanguage"
            :src="currentTheme.icons?.universalAccessSignLanguage"
          />
          <span
            v-text="
              $pgettext(
                'Accessibility menu: entry that links to sign-language resources or videos; feature name, Title Case.',
                'Sign Language'
              )
            "
          />
        </oc-button>
      </li>
    </oc-list>
  </oc-drop>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@ownclouders/web-pkg'
import { onMounted, unref, useTemplateRef } from 'vue'
import { OcDrop } from '@ownclouders/design-system/components'

const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const menu = useTemplateRef<InstanceType<typeof OcDrop>>('menu')

onMounted(() => {
  menu.value?.tippy?.setProps({
    onHidden: () => unref(menu).$el.focus(),
    onShown: () => unref(menu).$el.querySelector('a:first-of-type').focus()
  })
})
</script>

<style lang="scss" scoped>
.universal-access-dropdown-list {
  li {
    align-items: center;
    display: flex;
    margin: var(--oc-space-xsmall) 0;

    &:first-child {
      margin-top: 0;
    }

    &:last-child {
      margin-bottom: 0;
    }
  }

  a {
    gap: var(--oc-space-medium);
    justify-content: flex-start;
    min-height: 3rem;
    padding-left: var(--oc-space-small);
    width: 100%;

    &:focus,
    &:hover {
      background-color: var(--oc-color-background-hover);
      color: var(--oc-color-swatch-passive-default);
      text-decoration: none;
    }
  }
}
</style>
