<template>
  <div>
    <oc-button :id="buttonElementId" class="uk-width-expand">
      <span v-if="selectedOption">
        {{ selectedOption.displayValue }}
      </span>
      <span v-else>
        {{ setting.placeholder || $gettext('Please select') }}
      </span>
    </oc-button>
    <oc-drop
      :drop-id="dropElementId"
      :toggle="`#${buttonElementId}`"
      mode="click"
      close-on-click
      position="bottom-justify"
      :options="{ offset: 0, delayHide: 200, flip: false }"
      >
      <ul class="uk-list">
        <li
          v-for="(option, index) in setting.singleChoiceValue.options"
          :key="getOptionElementId(index)"
        >
          <label :for="getOptionElementId(index)">
            <input
              :id="getOptionElementId(index)"
              type="radio"
              class="oc-radiobutton"
              v-model="selectedOption"
              :value="option"
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
    dropElementId () {
      return `single-choice-drop-${this.bundle.identifier.bundleKey}-${this.setting.settingKey}`
    },
    buttonElementId () {
      return `single-choice-toggle-${this.bundle.identifier.bundleKey}-${this.setting.settingKey}`
    }
  },
  methods: {
    getOptionElementId (index) {
      return `${this.bundle.identifier.bundleKey}-${this.setting.settingKey}-${index}`
    },
    async onSelectedOption () {
      const value = {}
      if (this.selectedOption) {
        if (!isNil(this.selectedOption.intValue)) {
          value.intListValue = {
            value: [this.selectedOption ? this.selectedOption.intValue : null]
          }
        } else {
          value.stringListValue = {
            value: [this.selectedOption ? this.selectedOption.stringValue : null]
          }
        }
      }
      await this.$emit('onSave', {
        bundle: this.bundle,
        setting: this.setting,
        value
      })
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted () {
    if (!isNil(this.persistedValue)) {
      if (!isNil(this.persistedValue.intListValue)) {
        const selected = this.persistedValue.intListValue.value[0]
        const filtered = this.setting.singleChoiceValue.options.filter(option => option.intValue === selected)
        if (filtered.length > 0) {
          this.selectedOption = filtered[0]
        }
      } else {
        const selected = this.persistedValue.stringListValue.value[0]
        const filtered = this.setting.singleChoiceValue.options.filter(option => option.stringValue === selected)
        if (filtered.length > 0) {
          this.selectedOption = filtered[0]
        }
      }
    }
    // if not set, yet, apply default from settings bundle definition
    if (isNil(this.selectedOption)) {
      const defaults = this.setting.singleChoiceValue.options.filter(option => option.default)
      if (defaults.length === 1) {
        this.selectedOption = defaults[0]
      }
    }
  }
}
</script>
