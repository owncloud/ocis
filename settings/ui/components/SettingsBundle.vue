<template>
  <div class="uk-width-1-1 uk-width-2-3@m uk-width-1-2@l uk-width-1-3@xl">
    <h2 class="oc-mb-s">
      <translate>{{ bundle.displayName }}</translate>
    </h2>
    <oc-grid gutter="small">
      <template>
        <div class="uk-width-1-1" v-for="setting in bundle.settings" :key="setting.id">
          <label class="oc-label" :for="setting.id">{{ setting.displayName }}</label>
          <div class="uk-position-relative"
               :is="getSettingComponent(setting)"
               :id="setting.id"
               :bundle="bundle"
               :setting="setting"
               :persisted-value="getValue(setting)"
               @onSave="onSaveValue"
          />
        </div>
      </template>
    </oc-grid>
  </div>
</template>

<script>
import assign from 'lodash-es/assign'
import { mapGetters, mapActions } from 'vuex'
import SettingBoolean from './settings/SettingBoolean.vue'
import SettingMultiChoice from './settings/SettingMultiChoice.vue'
import SettingNumber from './settings/SettingNumber.vue'
import SettingSingleChoice from './settings/SettingSingleChoice.vue'
import SettingString from './settings/SettingString.vue'
import SettingUnknown from './settings/SettingUnknown.vue'
export default {
  name: 'SettingsBundle',
  props: {
    bundle: {
      type: Object,
      required: true
    }
  },
  computed: mapGetters(['getSettingsValue']),
  methods: {
    ...mapActions('Settings', ['saveValue']),
    getSettingComponent (setting) {
      return 'Setting' + setting.type[0].toUpperCase() + setting.type.substr(1)
    },
    getValue (setting) {
      return this.getSettingsValue({ settingId: setting.id })
    },
    async onSaveValue ({ bundle, setting, payload }) {
      payload = assign({}, payload, {
        bundleId: bundle.id,
        settingId: setting.id,
        accountUuid: 'me',
        resource: setting.resource
      })
      await this.saveValue({
        bundle,
        setting,
        payload
      })
    }
  },
  components: {
    SettingBoolean,
    SettingMultiChoice,
    SettingNumber,
    SettingSingleChoice,
    SettingString,
    SettingUnknown
  }
}
</script>
