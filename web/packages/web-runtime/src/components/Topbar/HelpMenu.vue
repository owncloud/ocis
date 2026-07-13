<template>
  <template v-if="isHelpMenuEnabled">
    <oc-button
      id="_helpMenuButton"
      ref="menuButton"
      v-oc-tooltip="
        $pgettext(
          'Top bar: tooltip for the “Help” icon; imperative sentence that opens a menu with license/help links.',
          'Open help menu'
        )
      "
      appearance="raw-inverse"
      variation="brand"
      :aria-label="$pgettext('Top bar: label for the help menu trigger; Title Case.', 'Help')"
    >
      <oc-icon name="question" fill-type="line" />
    </oc-button>
    <oc-drop ref="menu" toggle="#_helpMenuButton" mode="click" close-on-click padding-size="small">
      <oc-list class="help-menu-list">
        <li v-if="softwareLicenseUrl">
          <oc-button
            appearance="raw"
            type="a"
            :href="softwareLicenseUrl"
            target="_blank"
            rel="noopener noreferrer"
            data-testid="help-menu-software-license-link"
          >
            <oc-icon name="scales" fill-type="line" />
            <span
              v-text="
                $pgettext(
                  'Help menu: link label; opens the software license information page.',
                  'Software License Information'
                )
              "
            />
          </oc-button>
        </li>
        <li v-if="helpPageUrl">
          <oc-button
            appearance="raw"
            type="a"
            :href="helpPageUrl"
            target="_blank"
            rel="noopener noreferrer"
            data-testid="help-menu-help-page-link"
          >
            <oc-icon name="question" fill-type="line" />
            <span
              v-text="$pgettext('Help menu: link label; opens the help pages.', 'Help Pages')"
            />
          </oc-button>
        </li>
      </oc-list>
    </oc-drop>
  </template>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { computed, unref, onMounted, useTemplateRef } from 'vue'
import { useThemeStore } from '@ownclouders/web-pkg'
import { OcDrop } from '@ownclouders/design-system/components'
import { useGettext } from 'vue3-gettext'

const { $pgettext } = useGettext()
const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const softwareLicenseUrl = computed(() => unref(currentTheme).common?.urls?.softwareLicense)
const helpPageUrl = computed(() => unref(currentTheme).common?.urls?.helpPage)
const isHelpMenuEnabled = computed(() => !!(unref(softwareLicenseUrl) || unref(helpPageUrl)))

const menu = useTemplateRef<InstanceType<typeof OcDrop>>('menu')

onMounted(() => {
  menu.value?.tippy?.setProps({
    onHidden: () => unref(menu).$el.focus(),
    onShown: () => unref(menu).$el.querySelector('a:first-of-type').focus()
  })
})
</script>

<style lang="scss" scoped>
.help-menu-list {
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
