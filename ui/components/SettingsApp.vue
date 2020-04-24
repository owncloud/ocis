<template>
  <div>
    <oc-loader v-if="loading" />
    <div v-else class="uk-width-3-4@m uk-container uk-padding">
      <div class="uk-flex uk-flex-between uk-flex-middle">
        <h1 v-translate class="oc-page-title">Settings</h1>
      </div>
      <hr />
      <oc-alert v-if="extensions.length === 0" variation="primary" no-close>
        <p class="uk-flex uk-flex-middle">
          <oc-icon name="info" class="uk-margin-xsmall-right" />
          <translate>No settings available</translate>
        </p>
      </oc-alert>
      <template v-else>
        <settings-bundle
          v-for="bundle in visibleSettingsBundles"
          :key="'bundle-' + bundle.key"
          :bundle="bundle"
          class="uk-margin-top"
        />
      </template>
    </div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import SettingsBundle from "./SettingsBundle.vue";
export default {
  name: 'SettingsApp',
  components: {SettingsBundle},
  data () {
    return {
      loading: true,
      selectedExtension: undefined
    }
  },
  computed: {
    ...mapGetters('Settings', [
      'extensions',
      'getSettingsBundlesByExtension'
    ]),
    visibleSettingsBundles() {
      if (this.selectedExtension) {
        return this.getSettingsBundlesByExtension(this.selectedExtension)
      }
      return []
    }
  },
  methods: {
    ...mapActions('Settings', ['fetchSettingsBundles'])
  },
  async created () {
    await this.fetchSettingsBundles()
    if (this.extensions.length > 0) {
      this.selectedExtension = this.extensions[0]
    }
    this.loading = false;
  }
}
</script>
