<template>
  <oc-select
      v-model="selectedOptions"
      :clearable="false"
      :options="displayOptions"
      @input="onSelectedOption"
      multiple
  />
</template>

<script>
import isNil from 'lodash-es/isNil'
export default {
  name: 'SettingMultiChoice',
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
      selectedOptions: []
    }
  },
  computed: {
    displayOptions () {
      return this.setting.multiChoiceValue.options.map(val => val.displayValue)
    }
  },
  methods: {
    async onSelectedOption () {
      const values = []
      if (!isNil(this.selectedOptions)) {
        this.selectedOptions.forEach(displayValue => {
          const option = this.setting.multiChoiceValue.options.find(val => val.displayValue === displayValue)

          if (option.value.intValue) {
            values.push({ intValue: option.value.intValue })
          }
          if (option.value.stringValue) {
            values.push({ stringValue: option.value.stringValue })
          }
        })
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
      const selectedValues = []
      if (this.persistedValue.listValue.values) {
        this.persistedValue.listValue.values.forEach(value => {
          if (value.intValue) {
            selectedValues.push(value.intValue)
          }
          if (value.stringValue) {
            selectedValues.push(value.stringValue)
          }
        })
      }
      if (selectedValues.length === 0) {
        this.selectedOptions = []
      } else {
        this.selectedOptions = this.setting.multiChoiceValue.options.filter(option => {
          if (option.value.intValue) {
            return selectedValues.includes(option.value.intValue)
          }
          if (option.value.stringValue) {
            return selectedValues.includes(option.value.stringValue)
          }
          return false
        }).map(val => val.displayValue)
      }
    }
    // TODO: load the settings value of the authenticated user and set it in `selectedOptions`
    // if not set, yet, apply defaults from settings bundle definition
    if (this.selectedOptions === null) {
      this.selectedOptions = this.setting.multiChoiceValue.options
        .filter(option => option.default)
        .map(val => val.displayValue)
    }
  }
}
</script>
