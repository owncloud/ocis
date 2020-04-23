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
      close-on-click
      :drop-id="dropElementId"
      :toggle="`#${buttonElementId}`"
      mode="click"
      :options="{ offset: 0, delayHide: 0, flip: false }"
      >
      <ul class="uk-list">
        <li v-for="(option, index) in setting.singleChoiceValue.options" :key="`${setting.key}-${index}`">
          <oc-radio :label="option.displayValue" :value="selectedOption" @input="setSelectedOption(option)" />
        </li>
      </ul>
    </oc-drop>
  </div>
</template>

<script>
export default {
  name: 'SettingSingleChoice',
  props: {
    setting: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      selectedOption: null
    }
  },
  computed: {
    dropElementId() {
      return `single-choice-drop-${this.setting.key}`
    },
    buttonElementId() {
      return `single-choice-toggle-${this.setting.key}`
    },
  },
  methods: {
    setSelectedOption(option) {
      this.selectedOption = option
      // TODO: propagate selection to parent
      // TODO: show a spinner while the request for saving the value is running!
    }
  },
  mounted() {
    // TODO: load the settings value of the authenticated user and set it in `selected`
  }
}
</script>
