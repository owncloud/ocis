import { Ref, ref, unref } from 'vue'
import { SortDir, SortOptions, useSort } from '../../../../src/composables'
import { determineResourceTableSortFields } from '../../../../src/helpers/ui/resourceTable'
import {
  IncomingShareResource,
  OutgoingShareResource,
  Resource,
  ShareTypes
} from '@ownclouders/web-client'
import { getComposableWrapper } from '../../../../src/testing'
import { mock } from 'vitest-mock-extended'
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

  describe('sorting incoming shares by sharer', () => {
    const share = (id: string, name: string, sharers: string[]) =>
      mock<IncomingShareResource>({
        id,
        name,
        sdate: '2026-10-01',
        sharedBy: sharers.map((displayName) => ({ id: displayName, displayName }))
      })

    const sortBySharedBy = (shares: IncomingShareResource[], sortDir: Ref<SortDir>) => {
      const { items } = useSort<IncomingShareResource>({
        items: ref(shares),
        fields: determineResourceTableSortFields(shares[0]),
        sortBy: ref('sharedBy'),
        sortDir
      })
      return () => unref(items).map((i) => i.name)
    }

    it('sorts shares alphabetically by who shared them', () => {
      getComposableWrapper(() => {
        const sortDir = ref(SortDir.Asc)
        const names = sortBySharedBy(
          [
            share('1', 'budget.xlsx', ['Zoe Wagner']),
            share('2', 'letter.docx', ['Anna Bauer']),
            share('3', 'maths', ['Bert Klein'])
          ],
          sortDir
        )

        expect(names()).toEqual(['letter.docx', 'maths', 'budget.xlsx'])

        sortDir.value = SortDir.Desc
        expect(names()).toEqual(['budget.xlsx', 'maths', 'letter.docx'])
      })
    })

    it('places a share from several people by all sharers in alphabetical order', () => {
      getComposableWrapper(() => {
        const trip = share('1', 'trip', ['Zoe Wagner', 'Anna Bauer'])
        const names = sortBySharedBy(
          [trip, share('2', 'maths', ['Bert Klein']), share('3', 'letter.docx', ['Anna Bauer'])],
          ref(SortDir.Asc)
        )

        expect(names()).toEqual(['letter.docx', 'trip', 'maths'])
        // the avatars keep the order they were displayed in
        expect(trip.sharedBy.map(({ displayName }) => displayName)).toEqual([
          'Zoe Wagner',
          'Anna Bauer'
        ])
      })
    })

    it('places a share from a sharer without a display name first in A→Z', () => {
      getComposableWrapper(() => {
        const noName = mock<IncomingShareResource>({
          id: '1',
          name: 'noname.docx',
          sdate: '2026-10-01',
          sharedBy: [{ id: 'u1', displayName: undefined }]
        })
        const names = sortBySharedBy(
          [share('2', 'letter.docx', ['Anna Bauer']), noName],
          ref(SortDir.Asc)
        )

        expect(names()).toEqual(['noname.docx', 'letter.docx'])
      })
    })

    it('keeps shares from the same sharer in their existing relative order, in both directions', () => {
      getComposableWrapper(() => {
        const first = share('1', 'first.docx', ['Anna Bauer'])
        const second = share('2', 'second.docx', ['Anna Bauer'])
        const sortDir = ref(SortDir.Asc)
        const names = sortBySharedBy([first, second], sortDir)

        expect(names()).toEqual(['first.docx', 'second.docx'])

        sortDir.value = SortDir.Desc
        expect(names()).toEqual(['first.docx', 'second.docx'])
      })
    })

    it('orders sharer names inside the key the same way useSort compares rows (accent/case-insensitive, numeric)', () => {
      getComposableWrapper(() => {
        // Share from both "Lehrer 9" and "Lehrer 10": numeric-aware ordering puts
        // "Lehrer 9" before "Lehrer 10" inside the key, which changes how the whole
        // key compares against a share from "Lehrer 10" alone.
        const names = sortBySharedBy(
          [
            share('1', 'fromBoth.docx', ['Lehrer 9', 'Lehrer 10']),
            share('2', 'fromTen.docx', ['Lehrer 10'])
          ],
          ref(SortDir.Asc)
        )

        expect(names()).toEqual(['fromBoth.docx', 'fromTen.docx'])
      })
    })
  })

  describe('sorting outgoing shares by sharee', () => {
    it('does not reorder the sharees of a share', () => {
      getComposableWrapper(() => {
        const user = (displayName: string) => ({
          id: displayName,
          displayName,
          shareType: ShareTypes.user.value
        })
        const report = mock<OutgoingShareResource>({
          id: '1',
          name: 'report.pdf',
          sharedWith: [user('Richard Feynman'), user('Marie Curie'), user('Albert Einstein')]
        })
        const { items } = useSort<OutgoingShareResource>({
          items: ref([
            report,
            mock<OutgoingShareResource>({ id: '2', name: 'notes.md', sharedWith: [user('Bob')] })
          ]),
          fields: determineResourceTableSortFields(report),
          sortBy: ref('sharedWith'),
          sortDir: ref(SortDir.Asc)
        })

        expect(unref(items).map((i) => i.name)).toEqual(['report.pdf', 'notes.md'])
        // the avatars keep the order they were displayed in
        expect(report.sharedWith.map(({ displayName }) => displayName)).toEqual([
          'Richard Feynman',
          'Marie Curie',
          'Albert Einstein'
        ])
      })
    })
  })
})
