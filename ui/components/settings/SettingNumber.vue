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
        @keydown.esc="cancel"
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
import isNil from 'lodash/isNil'
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
    },
    persistedValue: {
      type: Object,
      required: false
    }
  },
  data () {
    return {
      initialValue: null,
      value: null
    }
  },
  computed: {
    isChanged () {
      return this.initialValue !== this.value
    },
    inputAttributes () {
      const attributes = {}
      if (!isNil(this.setting.intValue.min)) {
        attributes.min = this.setting.intValue.min
      }
      if (!isNil(this.setting.intValue.max)) {
        attributes.max = this.setting.intValue.max
      }
      if (!isNil(this.setting.intValue.step)) {
        attributes.step = this.setting.intValue.step
      }
      return attributes
    }
  },
  methods: {
    cancel () {
      this.value = this.initialValue
    },
    async applyValue () {
      const value = {
        intValue: this.value
      }
      await this.$emit('onSave', {
        bundle: this.bundle,
        setting: this.setting,
        value
      })
      // TODO: show a spinner while the request for saving the value is running!
      this.initialValue = this.value
    }
  },
  mounted () {
    if (!isNil(this.persistedValue)) {
      this.value = this.persistedValue.intValue
    }
    if (isNil(this.value) && !isNil(this.setting.intValue.default)) {
      this.value = this.setting.intValue.default
    }
    this.initialValue = this.value
  }
}
</script>
