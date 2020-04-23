<template>
  <div>
    <oc-checkbox v-model="value" :label="setting.boolValue.label" />
  </div>
</template>

<script>
export default {
  name: 'SettingBoolean',
  props: {
    setting: {
      type: Object,
      required: true
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
    // TODO: load the settings value of the authenticated user and apply it to the value
    if (this.value === null) {
      this.value = this.setting.boolValue.default
    }
  }
}
</script>
