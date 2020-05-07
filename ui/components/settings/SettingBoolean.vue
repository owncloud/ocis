<template>
  <div>
    <oc-checkbox v-model="value" :label="setting.boolValue.label" />
  </div>
</template>

<script>
import isNil from "lodash/isNil"
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
  data() {
    return {
      initialValue: null,
      value: null
    }
  },
  computed: {
    isChanged() {
      return this.initialValue !== this.value
    }
  },
  methods: {
    applyValue() {
      // TODO: propagate value to parent
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted() {
    if (!isNil(this.persistedValue)) {
      this.value = this.persistedValue.boolValue
    }
    if (isNil(this.value) && !isNil(this.setting.boolValue.default)) {
      this.value = this.setting.boolValue.default
    }
    this.initialValue = this.value
  }
}
</script>
