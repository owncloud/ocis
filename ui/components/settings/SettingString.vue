<template>
  <oc-grid flex>
    <div class="uk-width-expand">
      <oc-text-input
        v-model="value"
        :placeholder="setting.stringValue.placeholder"
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
import isNil from "lodash/isNil"
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
      this.initialValue = this.value
    },
    cancel() {
      this.value = this.initialValue
    }
  },
  mounted() {
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
