import LoadingScreen from '../../../../src/components/AppTemplates/PartialViews/LoadingScreen.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

describe('The external app loading screen component', () => {
  test('displays a spinner and a paragraph', () => {
    const wrapper = mount(LoadingScreen, {
      global: {
        stubs: {
          OcSpinner: true
        },
        plugins: [...defaultPlugins()]
      }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })
})
