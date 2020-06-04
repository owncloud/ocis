<template>
  <div>
    <div v-if="initialized" class="uk-width-3-4@m uk-container uk-padding">
      <oc-alert v-if="extensions.length === 0" variation="primary" no-close>
        <p class="uk-flex uk-flex-middle">
          <oc-icon name="info" class="uk-margin-xsmall-right" />
          <translate>No settings available</translate>
        </p>
      </oc-alert>
      <template v-else>
        <template v-if="selectedExtensionName">
          <div class="uk-flex uk-flex-between uk-flex-middle">
            <h1 class="oc-page-title">
              {{ selectedExtensionName }}
            </h1>
          </div>
          <hr />
        </template>
        <template v-if="settingsValuesLoaded">
          <settings-bundle
            v-for="bundle in selectedSettingsBundles"
            :key="'bundle-' + bundle.identifier.bundleKey"
            :bundle="bundle"
            class="uk-margin-top"
          />
        </template>
        <div class="uk-margin-top" v-else>
          <oc-loader :aria-label="$gettext('Loading settings values')" />
          <oc-alert :aria-hidden="true" varition="primary" no-close>
            <p v-translate>Loading settings values...</p>
          </oc-alert>
        </div>
      </template>
    </div>
    <oc-loader v-else />
  </div>
</template>

<script>
import { mapActions, mapGetters, mapMutations } from 'vuex'
import SettingsBundle from './SettingsBundle.vue'
export default {
  name: 'SettingsApp',
  components: { SettingsBundle },
  data () {
    return {
      loading: true,
      selectedExtension: undefined
    }
  },
  computed: {
    ...mapGetters(['settingsValuesLoaded', 'getNavItems']),
    ...mapGetters('Settings', [
      'extensions',
      'initialized',
      'getSettingsBundlesByExtension'
    ]),
    extensionRouteParam () {
      return this.$route.params.extension
    },
    selectedExtensionName () {
      return this.getExtensionName(this.selectedExtension)
    },
    selectedSettingsBundles () {
      if (this.selectedExtension) {
        return this.getSettingsBundlesByExtension(this.selectedExtension)
      }
      return []
    }
  },
  methods: {
    ...mapActions('Settings', ['initialize']),
    ...mapMutations(['ADD_NAV_ITEM']),
    resetSelectedExtension () {
      if (this.extensions.length > 0) {
        if (this.extensionRouteParam && this.extensions.includes(this.extensionRouteParam)) {
          this.selectedExtension = this.extensionRouteParam
        } else {
          this.selectedExtension = this.extensions[0]
        }
      }
    },
    resetMenuItems () {
      this.extensions.forEach(extension => {
        /*
         * TODO:
         * a) set up a map with possible extensions and icons?
         * or b) let extensions register app info like displayName + icon?
         */
        const navItem = {
          name: this.getExtensionName(extension),
          iconMaterial: 'application',
          route: {
            name: 'settings',
            path: `/${extension}`
          }
        }
        this.ADD_NAV_ITEM({
          extension: 'settings',
          navItem
        })
      })
      console.log(this.getNavItems('settings'))
    },
    getExtensionName (extension) {
      switch (extension) {
        case 'ocis-accounts': return this.$gettext('Account')
        case 'ocis-hello': return this.$gettext('Hello')
        default: return extension
      }
    }
  },
  async created () {
    await this.initialize()
  },
  watch: {
    initialized () {
      this.resetMenuItems()
      this.resetSelectedExtension()
    },
    extensionRouteParam () {
      this.resetSelectedExtension()
    }
  }
}
</script>
