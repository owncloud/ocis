import { shallowMount, mount, defaultPlugins } from '@ownclouders/web-test-helpers'
import Table from './OcTable.vue'

const fields = [
  {
    name: 'id',
    title: 'Id'
  },
  {
    name: 'resource',
    title: 'Resource',
    headerType: 'slot',
    type: 'slot'
  },
  {
    name: 'doubled',
    title: 'Doubled',
    type: 'callback',
    callback: function (value: number) {
      return `Double of ${value} is ${value * 2}`
    }
  }
]

const data = [
  {
    id: '4b136c0a-5057-11eb-ac70-eba264112003',
    resource: 'hello-world.txt',
    icon: 'text',
    doubled: 2
  },
  {
    id: '8468c9f0-5057-11eb-924b-934c6fd827a2',
    resource: 'I am a folder',
    icon: 'folder',
    doubled: 6
  },
  {
    id: '9c4cf97e-5057-11eb-8044-b3d5df9caa21',
    resource: 'this is fine.png',
    icon: 'image',
    doubled: 12
  }
]

describe('OcTable', () => {
  it('displays all field types', () => {
    const wrapper = mount(Table, {
      global: { plugins: defaultPlugins() },
      props: {
        fields,
        data
      },
      slots: {
        resourceHeader: '<span class="slot-header">Hello world!</span>',
        resource: `
        <div class="slot">
          <span>
            Hello world!
          </span>
        </div>
        `
      }
    })

    expect(wrapper.html().indexOf('4b136c0a-5057-11eb-ac70-eba264112003')).toBeGreaterThan(-1)
    expect(wrapper.html().indexOf('Double of 2 is 4')).toBeGreaterThan(-1)
    expect(wrapper.findAll('.slot').length).toBe(data.length)
    expect(wrapper.findAll('.slot-header').length).toBe(1)
  })

  it('hides header', () => {
    const wrapper = shallowMount(Table, {
      global: { plugins: defaultPlugins() },
      props: {
        fields,
        data,
        hasHeader: false
      }
    })

    expect(wrapper.findAll('oc-thead-stub').length).toBe(0)
  })

  it('enables hover effect', () => {
    const wrapper = shallowMount(Table, {
      global: { plugins: defaultPlugins() },
      props: {
        fields,
        data,
        hover: true
      }
    })

    expect(wrapper.attributes('class')).toContain('oc-table-hover')
  })

  it('extracts field title', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields: [
          {
            name: 'resource-name'
          }
        ],
        data: [
          {
            id: 'documents',
            name: 'Documents'
          }
        ]
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })

    expect(wrapper.html().indexOf('resource-name')).toBeGreaterThan(-1)
  })

  it('extracts cell props', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields: [
          {
            name: 'name',
            title: 'Name'
          },
          {
            name: 'description',
            title: 'Description',
            alignH: 'right',
            alignV: 'top',
            width: 'expand'
          }
        ],
        data: [
          {
            id: 'size',
            name: 'Size',
            description: 'Size of the resource'
          }
        ]
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })

    expect(wrapper.html().indexOf('alignh="right"')).toBeGreaterThan(-1)
    expect(wrapper.html().indexOf('alignv="top"')).toBeGreaterThan(-1)
    expect(wrapper.html().indexOf('width="expand"')).toBeGreaterThan(-1)
  })

  it('adds sticky header', () => {
    const wrapper = shallowMount(Table, {
      global: { plugins: defaultPlugins() },
      props: {
        fields,
        data,
        sticky: true
      }
    })

    expect(wrapper.attributes('class')).toContain('oc-table-sticky')
  })

  it('highlights a row', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields,
        data,
        highlighted: '4b136c0a-5057-11eb-ac70-eba264112003'
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })

    expect(wrapper.findAll('.oc-table-highlighted').length).toEqual(1)
  })

  it('highlights multiple rows', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields,
        data,
        highlighted: [
          '4b136c0a-5057-11eb-ac70-eba264112003',
          '8468c9f0-5057-11eb-924b-934c6fd827a2'
        ]
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })

    expect(wrapper.findAll('.oc-table-highlighted').length).toEqual(2)
  })

  it('adds data-item-id for rows', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields,
        data,
        highlighted: []
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })
    expect(wrapper.html().indexOf('data-item-id')).toBeGreaterThan(-1)
  })

  it('accepts itemDomSelector closure', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields,
        data,
        highlighted: [],
        itemDomSelector: (item: { id: string }) => ['custom', item.id].join('-')
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })
    data.forEach((item) => {
      expect(wrapper.find(['.oc-tbody-tr-custom', item.id].join('-')).exists()).toBeTruthy()
    })
  })

  it('emits contextmenu-clicked event upon right click on table row', async () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields,
        data,
        highlighted: []
      },
      global: { renderStubDefaultSlot: true, stubs: { OcTr: false }, plugins: defaultPlugins() }
    })
    await wrapper.find('.oc-tbody-tr').trigger('contextmenu')
    expect(wrapper.emitted().contextmenuClicked.length).toBe(1)
  })

  it('enable dragDrop should enable draggable on rows', () => {
    const wrapper = shallowMount(Table, {
      props: {
        fields,
        data,
        highlighted: [],
        dragDrop: true
      },
      global: { renderStubDefaultSlot: true, plugins: defaultPlugins() }
    })
    expect(wrapper.html().indexOf('draggable')).toBeGreaterThan(-1)
  })
})
