<template>
  <div>
    <oc-button :id="buttonElementId" class="uk-width-expand">
      <span v-if="selectedOptions !== null && selectedOptions.length > 0">
        {{ selectedOptionsDisplayValues }}
      </span>
      <span v-else>
        {{ setting.placeholder || $gettext('Please select') }}
      </span>
    </oc-button>
    <oc-drop
      :drop-id="dropElementId"
      :toggle="`#${buttonElementId}`"
      mode="click"
      position="bottom-justify"
      :options="{ offset: 0, delayHide: 200, flip: false }"
      >
      <ul class="uk-list">
        <li
          v-for="(option, index) in setting.multiChoiceValue.options"
          :key="getOptionElementId(index)"
        >
          <label :for="getOptionElementId(index)">
            <input
              :id="getOptionElementId(index)"
              type="checkbox"
              class="oc-checkbox"
              :value="option"
              v-model="selectedOptions"
              @change="onSelectedOption"
            />
            {{ option.displayValue }}
          </label>
        </li>
      </ul>
    </oc-drop>
  </div>
</template>

<script>
import isNil from 'lodash/isNil'
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
      selectedOptions: null
    }
  },
  computed: {
    selectedOptionsDisplayValues () {
      return Array.from(this.selectedOptions).map(option => option.displayValue).join(', ')
    },
    dropElementId () {
      return `multi-choice-drop-${this.bundle.identifier.bundleKey}-${this.setting.settingKey}`
    },
    buttonElementId () {
      return `multi-choice-toggle-${this.bundle.identifier.bundleKey}-${this.setting.settingKey}`
    }
  },
  methods: {
    getOptionElementId (index) {
      return `${this.bundle.identifier.bundleKey}-${this.setting.settingKey}-${index}`
    },
    async onSelectedOption () {
      const values = []
      if (!isNil(this.selectedOptions)) {
        this.selectedOptions.forEach(option => {
          if (option.value.intValue) {
            values.push({ intValue: option.value.intValue })
          }
          if (option.value.stringValue) {
            values.push({ stringValue: option.value.stringValue })
          }
        })
      }
      await this.$emit('onSave', {
        bundle: this.bundle,
        setting: this.setting,
        value: {
          listValue: {
            values
          }
        }
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
        })
      }
    }
    // TODO: load the settings value of the authenticated user and set it in `selectedOptions`
    // if not set, yet, apply defaults from settings bundle definition
    if (this.selectedOptions === null) {
      this.selectedOptions = this.setting.multiChoiceValue.options.filter(option => option.default)
    }
  }
}
</script>
