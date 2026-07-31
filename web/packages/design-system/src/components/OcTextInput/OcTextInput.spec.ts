import { shallowMount, mount, defaultPlugins } from '@ownclouders/web-test-helpers'
import OcTextInput from './OcTextInput.vue'
import { PasswordPolicy } from '../../helpers'
import { mock } from 'vitest-mock-extended'

interface Props {
  id?: string
  type?: string | 'text' | 'number' | 'email' | 'password' | 'date'
  modelValue?: string
  selectionRange?: [number, number]
  clearButtonEnabled?: boolean
  clearButtonAccessibleLabel?: string
  defaultValue?: string
  disabled?: boolean
  label: string
  warningMessage?: string
  errorMessage?: string
  fixMessageLine?: boolean
  descriptionMessage?: string
  readOnly?: boolean
  passwordPolicy?: PasswordPolicy
  generatePasswordMethod?: (...args: unknown[]) => string
}

const defaultProps = {
  label: 'label'
}

Object.assign(navigator, {
  clipboard: {
    writeText: vi.fn(),
    readText: vi.fn()
  }
})

// @vitest-environment jsdom
describe('OcTextInput', () => {
  function getShallowWrapper(props = {}) {
    return shallowMount(OcTextInput, {
      props: {
        ...defaultProps,
        ...props
      },
      global: { plugins: [...defaultPlugins()] }
    })
  }

  function getMountedWrapper(
    options: { props: Props; attachTo?: HTMLElement } = { props: { label: defaultProps.label } },
    passwordPolicy = { active: false, pass: false }
  ) {
    const passwordPolicyMock = mock<PasswordPolicy>()
    const policyRules = [
      {
        code: 'minLength',
        message: '%{param1}+ letters',
        format: ['8'],
        verified: passwordPolicy.pass
      }
    ]
    passwordPolicyMock.rules = policyRules
    passwordPolicyMock.missing.mockReturnValue({
      rules: policyRules
    })
    passwordPolicyMock.check.mockReturnValue(passwordPolicy.pass)

    if (passwordPolicy.active) {
      options.props = {
        ...(options.props || { label: defaultProps.label }),
        passwordPolicy: passwordPolicyMock
      }
    }

    const wrapper = mount(OcTextInput, {
      ...options,
      props: {
        'onUpdate:modelValue': (e) => wrapper.setProps({ modelValue: e }),
        modelValue: '',
        ...options.props
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })

    return wrapper
  }

  const selectors = {
    textInputMessage: '.oc-text-input-message span',
    clearInputButton: '.oc-text-input-btn-clear',
    inputField: '.oc-text-input',
    infoIcon: '.oc-text-input-message .oc-icon',
    showPasswordToggleBtn: '.oc-text-input-show-password-toggle',
    copyPasswordBtn: '.oc-text-input-copy-password-button',
    generatePasswordBtn: '.oc-text-input-generate-password-button'
  }

  describe('id prop', () => {
    const wrapper = getShallowWrapper({ id: 'test-input-id', descriptionMessage: 'hello' })
    it('should set provided id to the input', () => {
      expect(wrapper.find('input').attributes().id).toBe('test-input-id')
    })
    it('should set label target for provided id', () => {
      expect(wrapper.find('label').attributes().for).toBe('test-input-id')
    })
    it('should set message id according to provided id', () => {
      expect(wrapper.find(selectors.textInputMessage).attributes().id).toBe('test-input-id-message')
    })
  })

  describe('label prop', () => {
    it('should set provided label to the input', () => {
      const wrapper = getShallowWrapper()
      expect(wrapper.find('label').text()).toBe('label')
    })
  })

  describe('password input field', () => {
    describe('copy password button', () => {
      it('should not exist if type is not "password" or no value entered', () => {
        const wrapper = getMountedWrapper()
        expect(wrapper.find(selectors.copyPasswordBtn).exists()).toBeFalsy()
        const wrapper2 = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        expect(wrapper2.find(selectors.copyPasswordBtn).exists()).toBeFalsy()
      })
      it('should not exist if the input is disabled', () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', disabled: true, label: defaultProps.label }
        })
        expect(wrapper.find(selectors.copyPasswordBtn).exists()).toBeFalsy()
      })
      it('should exist if type is "password" and value entered', async () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        await wrapper.find(selectors.inputField).setValue('password')
        expect(wrapper.find(selectors.copyPasswordBtn).exists()).toBeTruthy()
      })
      it('should copy password to clipboard if clicked', async () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        await wrapper.find(selectors.inputField).setValue('password')
        await wrapper.find(selectors.copyPasswordBtn).trigger('click')
        expect(navigator.clipboard.writeText).toHaveBeenCalledWith('password')
      })
    })
    describe('show hide password toggle button', () => {
      it('should not exist if type is not "password" or no value entered', () => {
        const wrapper = getMountedWrapper()
        expect(wrapper.find(selectors.showPasswordToggleBtn).exists()).toBeFalsy()

        const wrapper2 = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        expect(wrapper2.find(selectors.showPasswordToggleBtn).exists()).toBeFalsy()
      })
      it('should not exist if the input is disabled', () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', disabled: true, label: defaultProps.label }
        })
        expect(wrapper.find(selectors.showPasswordToggleBtn).exists()).toBeFalsy()
      })
      it('should exist if type is "password" and value entered', async () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        await wrapper.find(selectors.inputField).setValue('password')
        expect(wrapper.find(selectors.showPasswordToggleBtn).exists()).toBeTruthy()
      })
      it('should show password plaintext/veiled if clicked', async () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        await wrapper.find(selectors.inputField).setValue('password')
        await wrapper.find(selectors.showPasswordToggleBtn).trigger('click')
        expect(wrapper.find(selectors.inputField).attributes().type).toBe('text')
        await wrapper.find(selectors.showPasswordToggleBtn).trigger('click')
        expect(wrapper.find(selectors.inputField).attributes().type).toBe('password')
      })
    })
    describe('generate password button', () => {
      it('should not exist if type is not "password" or prop "generatePasswordMethod" is not provided', () => {
        const wrapper = getMountedWrapper()
        expect(wrapper.find(selectors.generatePasswordBtn).exists()).toBeFalsy()

        const wrapper2 = getMountedWrapper({
          props: { type: 'password', label: defaultProps.label }
        })
        expect(wrapper2.find(selectors.generatePasswordBtn).exists()).toBeFalsy()
      })
      it('should not exist if the input is disabled', () => {
        const wrapper = getMountedWrapper({
          props: { type: 'password', disabled: true, label: defaultProps.label }
        })
        expect(wrapper.find(selectors.generatePasswordBtn).exists()).toBeFalsy()
      })
      it('should exist if type is "password" and prop "generatePasswordMethod" is provided', () => {
        const wrapper = getMountedWrapper({
          props: { generatePasswordMethod: vi.fn(), type: 'password', label: defaultProps.label }
        })
        expect(wrapper.find(selectors.generatePasswordBtn).exists()).toBeTruthy()
      })
      it('should fill input with generated password if clicked', async () => {
        const wrapper = getMountedWrapper({
          props: {
            generatePasswordMethod: vi.fn(() => 'PAssword12#!'),
            type: 'password',
            label: defaultProps.label
          }
        })
        await wrapper.find(selectors.generatePasswordBtn).trigger('click')
        expect((wrapper.find(selectors.inputField).element as HTMLInputElement).value).toEqual(
          'PAssword12#!'
        )
      })
    })
    describe('password policy', () => {
      it('should emit "passwordChallengeFailed" if password does not match criteria', async () => {
        const wrapper = getMountedWrapper(
          {
            props: { type: 'password', label: defaultProps.label }
          },
          { active: true, pass: false }
        )
        await wrapper.find(selectors.inputField).setValue('pass')
        expect(wrapper.emitted('passwordChallengeCompleted')).toBeFalsy()
      })
      it('should emit "passwordChallengeCompleted" if password matches criteria', async () => {
        const wrapper = getMountedWrapper(
          {
            props: { type: 'password', label: defaultProps.label }
          },
          { active: true, pass: true }
        )
        await wrapper.find(selectors.inputField).setValue('password123')
        expect(wrapper.emitted('passwordChallengeCompleted')).toBeTruthy()
      })
      it('displays error state if password does not match criteria', async () => {
        const wrapper = getMountedWrapper(
          {
            props: { type: 'password', label: defaultProps.label }
          },
          { active: true, pass: false }
        )
        await wrapper.find(selectors.inputField).setValue('pass')
        expect(wrapper.html()).toMatchSnapshot()
      })
      it('displays success state if password matches criteria', async () => {
        const wrapper = getMountedWrapper(
          {
            props: { type: 'password', label: defaultProps.label }
          },
          { active: true, pass: true }
        )
        await wrapper.find(selectors.inputField).setValue('password123')
        expect(wrapper.html()).toMatchSnapshot()
      })
    })
  })

  describe('when a description message is provided', () => {
    const wrapper = getShallowWrapper({ descriptionMessage: 'You should pass.' })
    it('should add the description class to the input message', () => {
      expect(wrapper.find(selectors.textInputMessage).attributes().class).toContain(
        'oc-text-input-description'
      )
    })
    it('should show the description message as the input message text', () => {
      expect(wrapper.find(selectors.textInputMessage).text()).toBe('You should pass.')
    })

    const mountedWrapper = getMountedWrapper({
      props: { descriptionMessage: 'You should pass.', label: defaultProps.label }
    })
    it('should show an info icon', () => {
      expect(mountedWrapper.find(selectors.infoIcon).exists()).toBe(true)
    })
  })

  describe('when a warning message is provided', () => {
    const wrapper = getShallowWrapper({ warningMessage: 'You may pass.' })
    it('should add the warning class to the input', () => {
      expect(wrapper.find('input').attributes().class).toContain('oc-text-input-warning')
    })
    it('should add the warning class to the input message', () => {
      expect(wrapper.find(selectors.textInputMessage).attributes().class).toContain(
        'oc-text-input-warning'
      )
    })
    it('should show the warning message as the input message text', () => {
      expect(wrapper.find(selectors.textInputMessage).text()).toBe('You may pass.')
    })
  })

  describe('when an error message is provided', () => {
    const wrapper = getShallowWrapper({ errorMessage: 'You shall not pass.' })
    it('should add the error class to the input', () => {
      expect(wrapper.find('input').attributes().class).toContain('oc-text-input-danger')
    })
    it('should add the error class to the input message', () => {
      expect(wrapper.find(selectors.textInputMessage).attributes().class).toContain(
        'oc-text-input-danger'
      )
    })
    it('should show the error message as the input message text', () => {
      expect(wrapper.find(selectors.textInputMessage).text()).toBe('You shall not pass.')
    })
    it('should set the input aria-invalid attribute to true', () => {
      expect(wrapper.find('input').attributes('aria-invalid')).toBe('true')
    })
  })

  describe('message priority', () => {
    it('should give error message top priority', () => {
      const wrapper = getShallowWrapper({
        errorMessage: 'You shall not pass.',
        warningMessage: 'You may pass.',
        descriptionMessage: 'Your should pass.'
      })
      const messageEl = wrapper.find('.oc-text-input-message span')
      expect(messageEl.attributes().class).toBe(
        'oc-text-input-description oc-text-input-warning oc-text-input-danger'
      )
      expect(messageEl.text()).toBe('You shall not pass.')
    })
    it('should give waring message priority over description message', () => {
      const wrapper = getShallowWrapper({
        warningMessage: 'You may pass.',
        descriptionMessage: 'Your should pass.'
      })
      const messageEl = wrapper.find(selectors.textInputMessage)
      expect(messageEl.attributes().class).toBe('oc-text-input-description oc-text-input-warning')
      expect(messageEl.text()).toBe('You may pass.')
    })
  })

  describe('type prop', () => {
    it.each(['text', 'number', 'email', 'password'])(
      'should set the provided type for the input',
      (type) => {
        const wrapper = getMountedWrapper({ props: { type: type, label: defaultProps.label } })
        expect(wrapper.find('input').attributes('type')).toBe(type)
      }
    )
  })

  describe('input events', () => {
    it('should emit an input event on typing', async () => {
      const wrapper = getShallowWrapper()
      expect(wrapper.emitted().input).toBeFalsy()
      await wrapper.find('input').setValue('a')
      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0][0]).toBe('a')
    })
  })

  describe('clear-button accessible label prop', () => {
    it('should set the aria label attribute if provided', () => {
      const wrapper = getShallowWrapper({
        clearButtonEnabled: true,
        modelValue: 'non-empty-value',
        clearButtonAccessibleLabel: 'test label'
      })
      expect(wrapper.find(selectors.clearInputButton).attributes('aria-label')).toBe('test label')
    })
  })

  describe('clear input', () => {
    it('has no clear button when it is disabled', () => {
      const wrapper = getShallowWrapper({
        modelValue: 'non-empty-value',
        clearButtonEnabled: true,
        disabled: true
      })
      expect(wrapper.find(selectors.clearInputButton).exists()).toBeFalsy()
    })

    it('has no clear button when it is enabled but the input is an empty string', () => {
      const wrapper = getShallowWrapper({
        clearButtonEnabled: true,
        modelValue: ''
      })

      expect(wrapper.find(selectors.clearInputButton).exists()).toBeFalsy()
    })

    it('has no clear button when it is enabled but the input is null', () => {
      const wrapper = getShallowWrapper({
        clearButtonEnabled: true,
        modelValue: null
      })

      expect(wrapper.find(selectors.clearInputButton).exists()).toBeFalsy()
    })

    it('has a clear button when it is enabled and the input contains content', () => {
      const wrapper = getShallowWrapper({
        clearButtonEnabled: true,
        modelValue: 'non-empty-value'
      })

      const btn = wrapper.find(selectors.clearInputButton)
      expect(btn.exists()).toBeTruthy()
    })

    it('clears the content on click', async () => {
      const wrapper = getMountedWrapper({
        props: {
          ...defaultProps,
          clearButtonEnabled: true,
          modelValue: 'non-empty-value'
        },
        attachTo: document.body
      })

      const btn = wrapper.find(selectors.clearInputButton)
      const input = wrapper.find(selectors.inputField)

      await btn.trigger('click')

      expect(wrapper.emitted('update:modelValue')[0][0]).toEqual(null)
      expect((input.element as HTMLInputElement).value).toEqual('')
      expect(document.activeElement.id).toBe(input.element.id)
    })
  })
})
