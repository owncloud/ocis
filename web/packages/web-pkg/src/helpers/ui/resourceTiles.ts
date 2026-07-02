import { Resource } from '@ownclouders/web-client'
import { SortDir, SortField } from '../../composables/sort'
import { Language } from 'vue3-gettext'

// just a dummy function to trick gettext tools
function $gettext(msg: string) {
  return msg
}

export const sortFields: SortField[] = [
  {
    label: $gettext('A-Z'),
    name: 'name',
    sortable: true,
    sortDir: SortDir.Asc
  },
  {
    label: $gettext('Z-A'),
    name: 'name',
    sortable: true,
    sortDir: SortDir.Desc
  },
  {
    label: $gettext('Newest'),
    name: 'mdate',
    sortable: (date: string) => new Date(date).valueOf(),
    sortDir: SortDir.Desc
  },
  {
    label: $gettext('Oldest'),
    name: 'mdate',
    sortable: (date: string) => new Date(date).valueOf(),
    sortDir: SortDir.Asc
  },
  {
    label: $gettext('Largest'),
    name: 'size',
    sortable: true,
    sortDir: SortDir.Desc
  },
  {
    label: $gettext('Smallest'),
    name: 'size',
    sortable: true,
    sortDir: SortDir.Asc
  },
  {
    label: $gettext('Remaining quota'),
    name: 'remainingQuota',
    prop: 'spaceQuota.remaining',
    sortable: true,
    sortDir: SortDir.Desc
  },
  {
    label: $gettext('Total quota'),
    name: 'totalQuota',
    prop: 'spaceQuota.total',
    sortable: true,
    sortDir: SortDir.Desc
  },
  {
    label: $gettext('Used quota'),
    name: 'usedQuota',
    prop: 'spaceQuota.used',
    sortable: true,
    sortDir: SortDir.Desc
  }
]

export const determineResourceTilesSortFields = (firstResource: Resource): SortField[] => {
  if (!firstResource) {
    return []
  }

  return sortFields.filter((field) =>
    Object.prototype.hasOwnProperty.call(firstResource, field.name)
  )
}

export const translateSortFields = (fields: SortField[], { $gettext }: Language): SortField[] => {
  return fields.map((field) => ({ ...field, label: $gettext(field.label) }))
}
