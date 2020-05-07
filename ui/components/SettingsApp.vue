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
        <settings-bundle
          v-for="bundle in selectedSettingsBundles"
          :key="'bundle-' + bundle.identifier.bundleKey"
          :bundle="bundle"
          class="uk-margin-top"
        />
      </template>
    </div>
    <oc-loader v-else />
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
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
    ...mapGetters('Settings', [
      'extensions',
      'initialized',
      'getSettingsBundlesByExtension'
    ]),
    extensionRouteParam () {
      return this.$route.params.extension
    },
    selectedExtensionName () {
      // TODO: extensions need to be registered with display names, separate from the settings bundles. until then: hardcoded translation
      if (this.selectedExtension === 'ocis-accounts') {
        return 'Account'
      } else if (this.selectedExtension === 'ocis-test') {
        return 'Test'
      }
      return this.selectedExtension
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
    resetSelectedExtension () {
      if (this.extensions.length > 0) {
        if (this.extensionRouteParam && this.extensions.includes(this.extensionRouteParam)) {
          this.selectedExtension = this.extensionRouteParam
        } else {
          this.selectedExtension = this.extensions[0]
        }
      }
    }
  },
  async created () {
    await this.initialize()
    this.resetSelectedExtension()
  },
  watch: {
    initialized () {
      this.resetSelectedExtension()
    },
    extensionRouteParam () {
      this.resetSelectedExtension()
    }
  }
}
</script>
