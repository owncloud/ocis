<template>
  <div>
    <oc-select
        v-model="selectedOption"
        :clearable="false"
        :options="displayOptions"
        @input="onSelectedOption"
       />
  </div>
</template>

<script>
import isNil from 'lodash/isNil'
export default {
  name: 'SettingSingleChoice',
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
      selectedOption: null
    }
  },
  computed: {
    displayOptions () {
      return this.setting.singleChoiceValue.options.map(val => val.displayValue)
    }
  },
  methods: {
    async onSelectedOption () {
      const values = []
      if (!isNil(this.selectedOption)) {
        const option = this.setting.singleChoiceValue.options.find(val => val.displayValue === this.selectedOption)

        if (option.value.intValue) {
          values.push({ intValue: option.value.intValue })
        }
        if (option.value.stringValue) {
          values.push({ stringValue: option.value.stringValue })
        }
      }
      const payload = {
        listValue: {
          values
        }
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
    if (!isNil(this.persistedValue) && !isNil(this.persistedValue.listValue)) {
      const selected = this.persistedValue.listValue.values[0]
      const filtered = this.setting.singleChoiceValue.options.filter(option => {
        if (selected.intValue) {
          return option.value.intValue === selected.intValue
        } else {
          return option.value.stringValue === selected.stringValue
        }
      })
      if (filtered.length > 0) {
        this.selectedOption = filtered[0].displayValue
      }
    }
    // if not set, yet, apply default from settings bundle definition
    if (isNil(this.selectedOption)) {
      const defaults = this.setting.singleChoiceValue.options.filter(option => option.default)
      if (defaults.length === 1) {
        this.selectedOption = defaults[0].displayValue
      }
    }
  }
}
</script>
