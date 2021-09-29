<template>
  <oc-grid flex>
    <div class="uk-width-expand">
      <oc-text-input
        v-model="value"
        :placeholder="setting.stringValue.placeholder"
        :label="setting.description"
        @keydown.enter="applyValue"
        @keydown.esc="cancel"
      />
    </div>
    <div v-if="isChanged">
      <oc-button @click="cancel" class="oc-ml-s">
        <translate>Cancel</translate>
      </oc-button>
      <oc-button @click="applyValue" class="oc-ml-s" variation="primary">
        <translate>Save</translate>
      </oc-button>
    </div>
  </oc-grid>
</template>

<script>
import isNil from 'lodash-es/isNil'
export default {
  name: 'SettingString',
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
    }
  },
  methods: {
    async applyValue () {
      const payload = {
        stringValue: this.value
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
      this.initialValue = this.value
    },
    cancel () {
      this.value = this.initialValue
    }
  },
  mounted () {
    if (!isNil(this.persistedValue)) {
      this.value = this.persistedValue.stringValue
    }
    if (isNil(this.value) && !isNil(this.setting.stringValue.default)) {
      this.value = this.setting.stringValue.default
    }
    this.initialValue = this.value
  }
}
</script>
