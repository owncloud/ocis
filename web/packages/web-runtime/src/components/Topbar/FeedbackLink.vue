<template>
  <div class="oc-flex">
    <oc-button
      v-oc-tooltip="descriptionOrFallback"
      type="a"
      :href="hrefOrFallback"
      target="_blank"
      appearance="raw-inverse"
      variation="brand"
      :aria-label="ariaLabelOrFallback"
      aria-describedby="oc-feedback-link-description"
    >
      <oc-icon name="feedback" />
    </oc-button>
    <p
      id="oc-feedback-link-description"
      class="oc-invisible-sr"
      aria-hidden="true"
      v-text="descriptionOrFallback"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'FeedbackLink',
  props: {
    href: {
      type: String,
      required: false,
      default: null
    },
    ariaLabel: {
      type: String,
      required: false,
      default: null
    },
    description: {
      type: String,
      required: false,
      default: null
    }
  },
  computed: {
    hrefOrFallback() {
      return this.href || 'https://owncloud.com/web-design-feedback'
    },
    ariaLabelOrFallback() {
      return this.ariaLabel || this.$gettext('ownCloud feedback survey')
    },
    descriptionOrFallback() {
      return (
        this.description ||
        this.$gettext(
          "Provide your feedback: We'd like to improve the web design and would be happy to hear your feedback. Thank you! Your ownCloud team."
        )
      )
    }
  }
})
</script>
