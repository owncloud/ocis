<template>
  <div
    class="oc-text-input-password-wrapper"
    :class="{
      'oc-text-input-password-wrapper-warning': hasWarning,
      'oc-text-input-password-wrapper-danger': hasError
    }"
  >
    <!-- Input label is handled in the OcTextInput component -->
    <!-- eslint-disable-next-line vuejs-accessibility/form-control-has-label -->
    <input
      v-bind="attrs"
      ref="passwordInput"
      :value="value"
      :type="showPassword ? 'text' : 'password'"
      :disabled="disabled"
      :aria-invalid="hasPasswordPolicyViolation"
      @input="$emit('input', $event)"
      @change="$emit('change', $event)"
    />
    <oc-button
      v-if="value && !disabled"
      v-oc-tooltip="showPassword ? $gettext('Hide password') : $gettext('Show password')"
      :aria-label="showPassword ? $gettext('Hide password') : $gettext('Show password')"
      class="oc-text-input-show-password-toggle oc-px-s oc-background-default"
      appearance="raw"
      size="small"
      @click="showPassword = !showPassword"
    >
      <oc-icon size="small" :name="showPassword ? 'eye-off' : 'eye'" />
    </oc-button>
    <oc-button
      v-if="value && !disabled"
      v-oc-tooltip="$gettext('Copy password')"
      :aria-label="$gettext('Copy password')"
      class="oc-text-input-copy-password-button oc-px-s oc-background-default"
      appearance="raw"
      size="small"
      @click="copyPasswordToClipboard"
    >
      <oc-icon size="small" :name="copyPasswordIcon" />
    </oc-button>
    <oc-button
      v-if="generatePasswordMethod && !disabled"
      v-oc-tooltip="$gettext('Generate password')"
      :aria-label="$gettext('Generate password')"
      class="oc-text-input-generate-password-button oc-px-s oc-background-default"
      appearance="raw"
      size="small"
      @click="generatePassword"
    >
      <oc-icon size="small" name="refresh" fill-type="line" />
    </oc-button>
  </div>
  <portal v-if="showPasswordPolicyInformation" to="app.design-system.password-policy">
    <div
      :id="passwordPolicyId"
      class="oc-flex oc-text-small oc-text-input-password-policy-rule-wrapper oc-pt-s"
      role="status"
      aria-live="polite"
      aria-atomic="true"
    >
      <span class="oc-invisible-sr" v-text="passwordPolicySummary" />
      <div
        v-for="(testedRule, index) in testedPasswordPolicy.rules"
        :key="index"
        class="oc-flex oc-flex-middle oc-text-input-password-policy-rule"
      >
        <oc-icon
          size="small"
          class="oc-mr-xs"
          :name="testedRule.verified ? 'checkbox-circle' : 'close-circle'"
          :variation="testedRule.verified ? 'success' : 'danger'"
          :accessible-label="
            testedRule.verified ? $gettext('Fulfilled') : $gettext('Not fulfilled')
          "
        />
        <span
          :class="[
            { 'oc-text-input-success': testedRule.verified },
            { 'oc-text-input-danger': !testedRule.verified }
          ]"
          v-text="getPasswordPolicyRuleMessage(testedRule)"
        ></span>
        <oc-contextual-helper
          v-if="testedRule.helperMessage"
          :text="testedRule.helperMessage"
          :title="$gettext('Password policy')"
        />
      </div>
    </div>
  </portal>
</template>

<script lang="ts" setup>
import { computed, ref, unref, watch, useAttrs } from 'vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcButton from '../OcButton/OcButton.vue'
import { useGettext } from 'vue3-gettext'
import { PasswordPolicy, PasswordPolicyRule } from '../../helpers'

interface Props {
  value: string
  passwordPolicy: PasswordPolicy
  generatePasswordMethod?: (...args: unknown[]) => string
  hasWarning: boolean
  hasError: boolean
  disabled: boolean
}

interface Emits {
  (e: 'passwordChallengeCompleted'): void
  (e: 'passwordChallengeFailed'): void
  (e: 'passwordGenerated', password: string): void
  (e: 'input', event: Event): void
  (e: 'change', event: Event): void
}

defineOptions({
  name: 'OCTextInputPassword',
  components: { OcButton, OcIcon },
  status: 'ready',
  release: '1.0.0'
})

const {
  value,
  passwordPolicy = {
    rules: [],
    check: () => false,
    missing: () => ({ rules: [] })
  },
  generatePasswordMethod = null,
  hasWarning = false,
  hasError = false,
  disabled = false
} = defineProps<Partial<Props>>()

const emit = defineEmits<Emits>()
const attrs = useAttrs()
const passwordInput = ref(null)
const { $gettext } = useGettext()
const showPassword = ref(false)
const copyPasswordIconInitial = 'file-copy'
const copyPasswordIcon = ref(copyPasswordIconInitial)

const inputId = computed(() => attrs.id as string)
const passwordPolicyId = computed(() => `${inputId.value}-password-policy`)

const showPasswordPolicyInformation = computed(() => {
  return !!Object.keys(passwordPolicy?.rules || {}).length
})

const testedPasswordPolicy = computed(() => {
  return passwordPolicy.missing(unref(value))
})

const hasPasswordPolicyViolation = computed(() => {
  if (!Object.keys(passwordPolicy?.rules || {})?.length) {
    return false
  }
  // Check if password has any policy violations or if there's an explicit error
  return !passwordPolicy.check(unref(value)) || hasError
})

const passwordPolicySummary = computed(() => {
  const tested = testedPasswordPolicy.value
  if (!tested || !tested.rules || tested.rules.length === 0) {
    return ''
  }

  const unmetRequirements = tested.rules.filter((rule) => !rule.verified)
  if (unmetRequirements.length === 0) {
    return $gettext('Password meets all requirements')
  }

  const messages = unmetRequirements.map((rule) => getPasswordPolicyRuleMessage(rule))
  return $gettext('Password requirements: %{requirements}', {
    requirements: messages.join(', ')
  })
})

const getPasswordPolicyRuleMessage = (rule: PasswordPolicyRule) => {
  const paramObj: Record<string, string> = {}

  for (let formatKey = 0; formatKey < rule.format.length; formatKey++) {
    paramObj[`param${formatKey + 1}`] = rule.format[formatKey]?.toString()
  }

  return $gettext(rule.message, paramObj)
}

const copyPasswordToClipboard = () => {
  navigator.clipboard.writeText(unref(value))
  copyPasswordIcon.value = 'check'
  setTimeout(() => (copyPasswordIcon.value = copyPasswordIconInitial), 500)
}

const generatePassword = () => {
  const generatedPassword = generatePasswordMethod()

  const inputEvent = new Event('input', { bubbles: true })
  Object.defineProperty(inputEvent, 'target', { value: passwordInput.value })

  emit('input', inputEvent)
  emit('passwordGenerated', generatedPassword)
}

watch(
  () => value,
  (value) => {
    if (!Object.keys(passwordPolicy?.rules || {})?.length) {
      return
    }

    if (!passwordPolicy.check(value)) {
      return emit('passwordChallengeFailed')
    }

    emit('passwordChallengeCompleted')
  }
)
</script>
<style lang="scss">
.oc-text-input-password {
  &-wrapper {
    display: flex;
    flex-direction: row;
    padding: 0;
    border-radius: 5px;
    border: 1px solid var(--oc-color-input-border);
    background-color: var(--oc-color-background-highlight);

    input {
      flex-grow: 2;
      border: none;

      &:focus {
        outline: none;
      }
    }

    &-warning,
    &-warning:focus {
      border-color: var(--oc-color-swatch-warning-default) !important;
      color: var(--oc-color-swatch-warning-default) !important;
    }

    &-danger,
    &-danger:focus {
      border-color: var(--oc-color-swatch-danger-default) !important;
      color: var(--oc-color-swatch-danger-default) !important;
    }

    &:focus-within {
      border-color: var(--oc-color-swatch-passive-default);
    }
  }

  &-policy-rule-wrapper {
    flex-direction: row;
    flex-wrap: wrap;
    column-gap: var(--oc-space-small);
  }
}
</style>
