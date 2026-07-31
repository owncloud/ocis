import { Resource, ShareResource } from '@ownclouders/web-client'
import { SortDir, SortField } from '../../composables/sort'

export const determineResourceTableSortFields = (firstResource: Resource): SortField[] => {
  if (!firstResource) {
    return []
  }

  return [
    {
      name: 'name',
      sortable: true,
      sortDir: SortDir.Asc
    },
    {
      name: 'size',
      sortable: true,
      sortDir: SortDir.Desc
    },
    {
      name: 'sharedWith',
      sortable: (sharedWith: ShareResource['sharedWith']) => {
        if (sharedWith.length > 0) {
          // Ensure the sharees are always sorted and that users
          // take precedence over groups. Then return a string with
          // all elements to ensure shares with multiple shares do
          // not appear mixed within others with a single one
          return sharedWith
            .sort((a, b) => {
              if (a.shareType !== b.shareType) {
                return a.shareType < b.shareType ? -1 : 1
              }
              return a.displayName < b.displayName ? -1 : 1
            })
            .map((e) => e.displayName)
            .join()
        }
        return false
      },
      sortDir: SortDir.Asc
    },
    {
      name: 'owner',
      sortable: 'displayName',
      sortDir: SortDir.Asc
    },
    {
      name: 'mdate',
      sortable: (date: string) => new Date(date).valueOf(),
      sortDir: SortDir.Desc
    },
    {
      name: 'sdate',
      sortable: (date: string) => new Date(date).valueOf(),
      sortDir: SortDir.Desc
    },
    {
      name: 'ddate',
      sortable: (date: string) => new Date(date).valueOf(),
      sortDir: SortDir.Desc
    }
  ].filter((field) => Object.prototype.hasOwnProperty.call(firstResource, field.name))
}
