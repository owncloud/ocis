import ErrorScreen from '../../../../src/components/AppTemplates/PartialViews/ErrorScreen.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

describe('The external app error screen component', () => {
  test('displays an icon and a paragraph', () => {
    const wrapper = mount(ErrorScreen, {
      props: {
        message: 'Error when loading the application'
      },
      global: {
        stubs: {
          OcIcon: true
        },
        plugins: [...defaultPlugins()]
      }
    })
    expect(wrapper.html()).toMatchSnapshot()
  })
})
