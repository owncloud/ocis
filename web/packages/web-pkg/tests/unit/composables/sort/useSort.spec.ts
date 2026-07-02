import { ref, unref } from 'vue'
import { SortDir, SortOptions, useSort } from '../../../../src/composables'
import { Resource } from '@ownclouders/web-client'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'

describe('useSort', () => {
  it('should be valid', () => {
    expect(useSort).toBeDefined()
  })

  it('does not sort if no sort field was given', () => {
    getComposableWrapper(() => {
      const input: SortOptions<any> = {
        items: [{ id: '3' }, { id: '4' }, { id: '6' }, { id: '1' }, { id: '2' }, { id: '5' }],
        fields: [],
        sortBy: ref(null),
        sortDir: ref(null),
        routeName: 'mocked'
      }

      const { items } = useSort(input)

      expect(unref(items)).toMatchObject([
        { id: '3' },
        { id: '4' },
        { id: '6' },
        { id: '1' },
        { id: '2' },
        { id: '5' }
      ])
    })
  })

  describe('sorting resources', () => {
    const resources: Resource[] = [
      { id: '1', name: 'c.png', path: '', webDavPath: '', mdate: '2', spaceId: '1' },
      { id: '2', name: 'Dir4', path: '', webDavPath: '', mdate: '4', type: 'folder', spaceId: '1' },
      { id: '3', name: 'a.png', path: '', webDavPath: '', mdate: '3', spaceId: '1' },
      { id: '4', name: 'A.png', path: '', webDavPath: '', mdate: '6', spaceId: '1' },
      { id: '5', name: 'dir2', path: '', webDavPath: '', mdate: '7', type: 'folder', spaceId: '1' },
      { id: '6', name: 'b.png', path: '', webDavPath: '', mdate: '1', spaceId: '1' },
      { id: '7', name: 'Dir1', path: '', webDavPath: '', mdate: '5', type: 'folder', spaceId: '1' },
      {
        id: '8',
        name: 'dir11',
        path: '',
        webDavPath: '',
        mdate: '8',
        type: 'folder',
        spaceId: '1'
      },
      { id: '9', name: 'dir3', path: '', webDavPath: '', mdate: '9', type: 'folder', spaceId: '1' }
    ]

    it('sorts resources by name', () => {
      getComposableWrapper(() => {
        const sortDir = ref(SortDir.Asc)
        const input = {
          items: ref<Resource[]>(resources),
          fields: [
            {
              name: 'name',
              sortable: true
            },
            {
              name: 'mdate',
              sortable: true
            }
          ],
          sortBy: ref('name'),
          sortDir
        }

        const { items } = useSort<Resource>(input)

        expect(unref(items).map((i) => i.name)).toMatchObject([
          'Dir1',
          'dir2',
          'dir3',
          'Dir4',
          'dir11',
          'a.png',
          'A.png',
          'b.png',
          'c.png'
        ])

        sortDir.value = SortDir.Desc
        expect(unref(items).map((i) => i.name)).toMatchObject([
          'c.png',
          'b.png',
          'a.png',
          'A.png',
          'dir11',
          'Dir4',
          'dir3',
          'dir2',
          'Dir1'
        ])
      })
    })
  })
})
