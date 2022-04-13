<template>
  <div>
    <oc-checkbox v-model="value" :label="setting.boolValue.label" @change="applyValue" />
  </div>
</template>

<script>
import isNil from 'lodash-es/isNil'
export default {
  name: 'SettingBoolean',
  props: {
    bundle: {
      type: Object,
      required: true
    },
    setting: {
      type: Object,
      required: true
    },
    persistedValue: {
      type: Object,
      required: false
    }
  },
  data () {
    return {
      value: null
    }
  },
  methods: {
    async applyValue () {
      const payload = {
        boolValue: this.value
      }
      if (!isNil(this.persistedValue)) {
        payload.id = this.persistedValue.id
      }
      await this.$emit('onSave', {
        bundle: this.bundle,
        setting: this.setting,
        payload
      })
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted () {
    if (!isNil(this.persistedValue)) {
      this.value = this.persistedValue.boolValue
    }
    if (isNil(this.value) && !isNil(this.setting.boolValue.default)) {
      this.value = this.setting.boolValue.default
    }
  }
}
</script>
