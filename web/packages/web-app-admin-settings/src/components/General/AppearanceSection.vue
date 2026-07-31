<template>
  <div>
    <h2 class="oc-py-s" v-text="$gettext('Appearance')" />
    <div>
      <div class="oc-flex oc-flex-middle">
        <h3 v-text="$gettext('Logo')" />
        <oc-button
          v-if="menuItems.length"
          :id="`logo-context-btn`"
          v-oc-tooltip="$gettext('Show context menu')"
          :aria-label="$gettext('Show context menu')"
          appearance="raw"
          class="oc-ml-s"
        >
          <oc-icon name="more-2" />
        </oc-button>
        <oc-drop
          :drop-id="`space-context-drop`"
          :toggle="`#logo-context-btn`"
          mode="click"
          close-on-click
          padding-size="small"
        >
          <context-action-menu :menu-sections="menuSections" :action-options="actionOptions" />
        </oc-drop>
      </div>
      <div>
        <div class="logo-wrapper">
          <img alt="" :src="currentTheme.logo.topbar" class="oc-p-xs" />
        </div>
        <label for="logo-upload-input" class="oc-invisible-sr">
          {{ $pgettext('Accesibility label to change logo', 'Change logo') }}
        </label>
        <input
          id="logo-upload-input"
          ref="logoInput"
          type="file"
          name="file"
          tabindex="-1"
          :accept="supportedLogoMimeTypesAcceptValue"
          @change="uploadImage"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref, VNodeRef, ref } from 'vue'
import { ContextActionMenu, useThemeStore } from '@ownclouders/web-pkg'
import {
  useGeneralActionsResetLogo,
  useGeneralActionsUploadLogo
} from '../../composables/actions/general'
import { supportedLogoMimeTypes } from '../../defaults'
import { storeToRefs } from 'pinia'

defineOptions({
  name: 'AppearanceSection'
})

const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const logoInput: VNodeRef = ref(null)

const { actions: resetLogoActions } = useGeneralActionsResetLogo()
const { actions: uploadLogoActions, uploadImage } = useGeneralActionsUploadLogo({
  imageInput: logoInput
})

const menuItems = computed(() =>
  [...unref(uploadLogoActions), ...unref(resetLogoActions)].filter((i) => i.isVisible())
)

const actionOptions = computed(() => ({
  resources: unref(menuItems)
}))

const menuSections = computed(() => [
  {
    name: 'primaryActions',
    items: unref(menuItems)
  }
])

const supportedLogoMimeTypesAcceptValue = supportedLogoMimeTypes.join(',')
</script>

<style lang="scss" scoped>
#logo-upload-input {
  display: none;
}
.logo-wrapper {
  width: 280px;
  min-width: 280px;
  aspect-ratio: 16/9;
  max-height: 158px;

  img {
    border-radius: 10px;
    width: 100%;
    max-height: 100%;
    object-fit: cover;
    background: var(--oc-color-swatch-brand-default);
  }
}
</style>
