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
      return `single-choice-drop-${this.setting.id}`
    },
    buttonElementId () {
      return `single-choice-toggle-${this.setting.id}`
    }
  },
  methods: {
    getOptionElementId (index) {
      return `${this.setting.id}-${index}`
    },
    async onSelectedOption () {
      const values = []
      if (!isNil(this.selectedOption)) {
        if (this.selectedOption.value.intValue) {
          values.push({ intValue: this.selectedOption.value.intValue })
        }
        if (this.selectedOption.value.stringValue) {
          values.push({ stringValue: this.selectedOption.value.stringValue })
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
        this.selectedOption = filtered[0]
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
