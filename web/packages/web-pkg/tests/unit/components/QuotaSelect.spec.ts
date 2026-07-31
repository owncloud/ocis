import QuotaSelect from '../../../src/components/QuotaSelect.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

describe('QuotaSelect', () => {
  describe('method "optionSelectable"', () => {
    it('should return true while option selectable property is not false', () => {
      const { wrapper } = getWrapper()
      expect(
        (wrapper.vm as any).optionSelectable({ value: 1, displayValue: '', selectable: true })
      ).toBeTruthy()
      expect((wrapper.vm as any).optionSelectable({ value: 1, displayValue: '' })).toBeTruthy()
    })
    it('should return false while option selectable property is false', () => {
      const { wrapper } = getWrapper()
      expect(
        (wrapper.vm as any).optionSelectable({ value: 1, displayValue: '', selectable: false })
      ).toBeFalsy()
    })
  })
  describe('method "createOption"', () => {
    it('should create option', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).createOption('3')).toEqual({
        displayValue: '3 GB',
        value: 3 * Math.pow(10, 9)
      })
    })
    it('should contain error property while maxQuota will be exceeded', () => {
      const { wrapper } = getWrapper({ maxQuota: 3 * Math.pow(10, 9) })
      expect((wrapper.vm as any).createOption('2000')).toHaveProperty('error')
    })
    it('should contain error property while creating an invalid option', () => {
      const { wrapper } = getWrapper()
      expect((wrapper.vm as any).createOption('lorem ipsum')).toHaveProperty('error')
      expect((wrapper.vm as any).createOption('1,')).toHaveProperty('error')
      expect((wrapper.vm as any).createOption('1.')).toHaveProperty('error')
    })
  })
  describe('method "setOptions"', () => {
    it('should set options to default options', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).setOptions()
      expect((wrapper.vm as any).options).toEqual((wrapper.vm as any).DEFAULT_OPTIONS)
    })
    it('should contain default options and user defined option if set', () => {
      const { wrapper } = getWrapper({ totalQuota: 45 * Math.pow(10, 9) })
      ;(wrapper.vm as any).setOptions()
      expect((wrapper.vm as any).options).toEqual(
        expect.arrayContaining([
          ...(wrapper.vm as any).DEFAULT_OPTIONS,
          {
            displayValue: '45 GB',
            value: 45 * Math.pow(10, 9),
            selectable: true
          }
        ])
      )
    })
    it('should only contain lower or equal options when max quota is set', () => {
      const { wrapper } = getWrapper({
        totalQuota: 2 * Math.pow(10, 9),
        maxQuota: 4 * Math.pow(10, 9)
      })
      ;(wrapper.vm as any).setOptions()
      expect((wrapper.vm as any).options).toEqual(
        expect.arrayContaining([
          {
            displayValue: '1 GB',
            value: Math.pow(10, 9)
          },
          {
            displayValue: '2 GB',
            value: 2 * Math.pow(10, 9)
          }
        ])
      )
    })
    it('should contain a non selectable option if preset quota is higher than max quota', () => {
      const { wrapper } = getWrapper({
        totalQuota: 100 * Math.pow(10, 9),
        maxQuota: 4 * Math.pow(10, 9)
      })
      ;(wrapper.vm as any).setOptions()
      expect((wrapper.vm as any).options).toEqual(
        expect.arrayContaining([
          {
            displayValue: '1 GB',
            value: Math.pow(10, 9)
          },
          {
            displayValue: '2 GB',
            value: 2 * Math.pow(10, 9)
          },
          {
            displayValue: '100 GB',
            value: 100 * Math.pow(10, 9),
            selectable: false
          }
        ])
      )
    })
  })
})

function getWrapper({ totalQuota = 10 * Math.pow(10, 9), maxQuota = 0 } = {}) {
  return {
    wrapper: shallowMount(QuotaSelect, {
      data: () => {
        return {
          selectedOption: {
            value: 10 * Math.pow(10, 9)
          },
          options: []
        }
      },
      props: {
        totalQuota,
        maxQuota,
        title: 'Personal quota'
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
