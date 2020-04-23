<template>
  <oc-grid flex>
    <div class="uk-width-expand">
      <oc-text-input
        type="number"
        v-model="value"
        v-bind="inputAttributes"
        :placeholder="setting.intValue.placeholder"
        :label="setting.description"
        @keydown.enter="applyValue"
      />
    </div>
    <div v-if="isChanged">
      <oc-button @click="cancel" class="uk-margin-xsmall-left">
        <translate>Cancel</translate>
      </oc-button>
      <oc-button @click="applyValue" class="uk-margin-xsmall-left" variation="primary">
        <translate>Save</translate>
      </oc-button>
    </div>
  </oc-grid>
</template>

<script>
export default {
  name: 'SettingNumber',
  props: {
    bundle: {
      type: Object,
      required: true
    },
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
    },
    inputAttributes() {
      const attributes = {}
      if (this.setting.intValue.min !== null && this.setting.intValue.min !== undefined) {
        attributes.min = this.setting.intValue.min
      }
      if (this.setting.intValue.max !== null && this.setting.intValue.max !== undefined) {
        attributes.max = this.setting.intValue.max
      }
      if (this.setting.intValue.step !== null && this.setting.intValue.step !== undefined) {
        attributes.step = this.setting.intValue.step
      }
      return attributes
    }
  },
  methods: {
    cancel() {
      this.value = this.initialValue
    },
    applyValue() {
      // TODO: propagate value to parent
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted() {
    // TODO: load the settings value of the authenticated user and apply it to the value
  }
}
</script>
