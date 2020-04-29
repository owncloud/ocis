<template>
  <div class="uk-width-1-1 uk-width-2-3@m uk-width-1-2@l">
    <div class="uk-text-bold uk-margin-small-bottom">
      <translate>{{ bundle.displayName }}</translate>
    </div>
    <oc-grid gutter="small">
      <div class="uk-width-1-1" v-for="setting in bundle.settings" :key="getElementId(bundle, setting)">
        <label class="oc-label" :for="getElementId(bundle, setting)">{{ setting.displayName }}</label>
        <div class="uk-position-relative"
             :is="getSettingComponent(setting)"
             :id="getElementId(bundle, setting)"
             :bundle="bundle"
             :setting="setting"
        />
      </div>
    </oc-grid>
  </div>
</template>

<script>
import SettingBoolean from "./settings/SettingBoolean.vue";
import SettingMultiChoice from "./settings/SettingMultiChoice.vue";
import SettingNumber from "./settings/SettingNumber.vue";
import SettingSingleChoice from "./settings/SettingSingleChoice.vue";
import SettingString from "./settings/SettingString.vue";
import SettingUnknown from "./settings/SettingUnknown.vue";

export default {
  name: 'SettingsBundle',
  props: {
    bundle: {
      type: Object,
      required: true
    }
  },
  methods: {
    getElementId(bundle, setting) {
      return `setting-${bundle.identifier.bundleKey}-${setting.settingKey}`
    },
    getSettingComponent(setting) {
      return 'Setting' + setting.type[0].toUpperCase() + setting.type.substr(1)
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
