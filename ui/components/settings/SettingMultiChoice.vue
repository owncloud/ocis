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
  name: 'SettingMultiChoice',
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
      selectedOptions: null
    }
  },
  computed: {
    selectedOptionsDisplayValues() {
      return Array.from(this.selectedOptions).map(option => option.displayValue).join(', ')
    },
    dropElementId() {
      return `multi-choice-drop-${this.bundle.key}-${this.setting.key}`
    },
    buttonElementId() {
      return `multi-choice-toggle-${this.bundle.key}-${this.setting.key}`
    },
  },
  methods: {
    getOptionElementId(index) {
      return `${this.bundle.key}-${this.setting.key}-${index}`
    },
    onSelectedOption() {
      // TODO: propagate selection to parent
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted() {
    this.selectedOptions = null
    // TODO: load the settings value of the authenticated user and set it in `selectedOptions`
    // if not set, yet, apply defaults from settings bundle definition
    if (this.selectedOptions === null) {
      this.selectedOptions = this.setting.multiChoiceValue.options.filter(option => option.default)
    }
  }
}
</script>
