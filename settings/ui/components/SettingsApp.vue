<template>
  <div class="oc-p">
    <div class="uk-flex uk-flex-column" id="settings-app">
      <template v-if="initialized">
        <oc-alert v-if="extensions.length === 0" variation="primary" no-close>
          <p class="uk-flex uk-flex-middle">
            <oc-icon name="info" class="oc-mr-s" />
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
              v-for="bundle in selectedBundles"
              :key="'bundle-' + bundle.id"
              :bundle="bundle"
              class="oc-mt"
            />
          </template>
          <div class="oc-mt" v-else>
            <oc-loader :aria-label="$gettext('Loading personal settings')" />
            <oc-alert :aria-hidden="true" varition="primary" no-close>
              <p v-translate>Loading personal settings...</p>
            </oc-alert>
          </div>
        </template>
      </template>
      <oc-loader v-else />
    </div>
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
      'getBundlesByExtension'
    ]),
    extensionRouteParam () {
      return this.$route.params.extension
    },
    selectedExtensionName () {
      return this.getExtensionName(this.selectedExtension)
    },
    selectedBundles () {
      if (this.selectedExtension) {
        return this.getBundlesByExtension(this.selectedExtension)
      }
      return []
    }
  },
  methods: {
    ...mapActions('Settings', ['initialize']),
    ...mapMutations(['ADD_NAV_ITEM']),
    resetSelectedExtension () {
      if (this.extensions.length > 0) {
        if (
          this.extensionRouteParam &&
          this.extensions.includes(this.extensionRouteParam)
        ) {
          this.selectedExtension = this.extensionRouteParam
        } else {
          this.selectedExtension = this.extensions[0]
        }
      }
    },
    resetMenuItems () {
      this.extensions.forEach((extension) => {
        /*
         * TODO:
         * a) set up a map with possible extensions and icons?
         * or b) let extensions register app info like displayName + icon?
         * https://github.com/owncloud/ocis/settings/issues/27
         */
        const navItem = {
          name: this.getExtensionName(extension),
          iconMaterial: this.getExtensionIcon(extension),
          route: {
            name: 'settings',
            path: `/settings/${extension}`
          },
          menu: 'user'
        }
        this.ADD_NAV_ITEM({
          extension: 'settings',
          navItem
        })
      })
    },
    getExtensionName (extension) {
      extension = extension || ''
      switch (extension) {
        case 'ocis-accounts':
          return this.$gettext('Account')
        case 'ocis-hello':
          return this.$gettext('Hello')
        default: {
          const shortenedName = extension.replace('ocis-', '')
          return shortenedName.charAt(0).toUpperCase() + shortenedName.slice(1)
        }
      }
    },
    getExtensionIcon (extension) {
      extension = extension || ''
      switch (extension) {
        case 'ocis-accounts':
          return 'account_circle'
        case 'ocis-hello':
          return 'tag_faces'
        default:
          return 'application'
      }
    }
  },
  created () {
    this.initialize()
  },
  watch: {
    '$language.current': {
      handler () {
        this.resetMenuItems()
      }
    },
    initialized: {
      handler () {
        this.resetMenuItems()
        this.resetSelectedExtension()
      },
      immediate: true
    },
    extensionRouteParam () {
      this.resetSelectedExtension()
    }
  }
}
</script>
