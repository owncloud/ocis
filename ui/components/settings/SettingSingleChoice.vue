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
              @input="onSelectedOption"
            />
            {{ option.displayValue }}
          </label>
        </li>
      </ul>
    </oc-drop>
  </div>
</template>

<script>
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
  data() {
    return {
      selectedOption: null
    }
  },
  computed: {
    dropElementId() {
      return `single-choice-drop-${this.bundle.identifier.bundleKey}-${this.setting.settingKey}`
    },
    buttonElementId() {
      return `single-choice-toggle-${this.bundle.identifier.bundleKey}-${this.setting.settingKey}`
    },
  },
  methods: {
    getOptionElementId(index) {
      return `${this.bundle.identifier.bundleKey}-${this.setting.settingKey}-${index}`
    },
    onSelectedOption() {
      // TODO: propagate selection to parent
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted() {
    this.selectedOption = null
    // TODO: load the settings value of the authenticated user and set it in `selectedOption`
    // if not set, yet, apply default from settings bundle definition
    if (this.selectedOption === null) {
      const defaults = this.setting.singleChoiceValue.options.filter(option => option.default)
      if (defaults.length === 1) {
        this.selectedOption = defaults[0]
      }
    }
  }
}
</script>
