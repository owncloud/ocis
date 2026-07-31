import { shallowMount } from '@ownclouders/web-test-helpers'
import Avatars from './OcAvatars.vue'

const items = [
  {
    id: 'bob',
    username: 'bob',
    displayName: 'Bob',
    avatar:
      'https://images.unsplash.com/photo-1610216705422-caa3fcb6d158?ixid=MXwxMjA3fDB8MHxzZWFyY2h8MTB8fGZhY2V8ZW58MHwyfDB8&ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=60',
    shareType: 0
  },
  {
    id: 'link0',
    link: 'fake url content',
    name: 'link 0',
    shareType: 3
  },
  {
    id: 'marie',
    username: 'marie',
    displayName: 'Marie',
    avatar:
      'https://images.unsplash.com/photo-1584308972272-9e4e7685e80f?ixid=MXwxMjA3fDB8MHxzZWFyY2h8Mzh8fGZhY2V8ZW58MHwyfDB8&ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=60',
    shareType: 0
  },
  {
    id: 'john',
    username: 'john',
    displayName: 'John Richards Emperor of long names',
    shareType: 0
  },
  {
    id: 'link1',
    link: 'fake url content',
    name: 'link 1',
    shareType: 3
  }
]

describe('OcAvatars', () => {
  it('displays tooltip', () => {
    const OcTooltip = vi.fn()
    const wrapper = shallowMount(Avatars, {
      props: {
        items,
        maxDisplayed: 2,
        isTooltipDisplayed: true,
        accessibleDescription: 'List of users'
      },
      global: {
        directives: {
          OcTooltip
        }
      }
    })

    expect(OcTooltip.mock.calls[0][1].value).toMatch('Bob, Marie +3')
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('prefers avatars over links when maxDisplayed is exceeded', () => {
    const OcTooltip = vi.fn()
    const wrapper = shallowMount(Avatars, {
      props: {
        items,
        maxDisplayed: 3,
        isTooltipDisplayed: true,
        accessibleDescription: 'List of users'
      },
      global: {
        directives: {
          OcTooltip
        }
      }
    })

    expect(OcTooltip.mock.calls[0][1].value).toMatch(
      'Bob, Marie, John Richards Emperor of long names +2'
    )
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('shows avatars first and links last', () => {
    const OcTooltip = vi.fn()
    const wrapper = shallowMount(Avatars, {
      props: {
        items,
        isTooltipDisplayed: true,
        accessibleDescription: 'List of users'
      },
      global: {
        directives: {
          OcTooltip
        }
      }
    })

    expect(OcTooltip.mock.calls[0][1].value).toMatch(
      'Bob, Marie, John Richards Emperor of long names, link 0, link 1'
    )
    expect(wrapper.html()).toMatchSnapshot()
  })
})
