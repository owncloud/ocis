import { IncomingShareResource } from '@ownclouders/web-client'
import { computed, Ref, unref } from 'vue'

export interface GroupingSettings {
  groupingBy: string
  showGroupingOptions: boolean
  groupingFunctions: {
    [key: string]: (row: IncomingShareResource) => string | void
  }
  sortGroups: {
    [key: string]: (groups: { name: string }[]) => { name: string }[]
  }
}

export const useGroupingSettings = ({
  sortBy,
  sortDir
}: {
  sortBy: Ref<string>
  sortDir: Ref<string>
}) => {
  const groupingSettings = computed(() => {
    return {
      groupingBy: localStorage.getItem('grouping-shared-with-me') || 'Shared on',
      showGroupingOptions: true,
      groupingFunctions: {
        'Name alphabetically': function (row: IncomingShareResource) {
          localStorage.setItem('grouping-shared-with-me', 'Name alphabetically')
          if (!isNaN(Number(row.name.charAt(0)))) {
            return '#'
          }
          if (row.name.charAt(0) === '.') {
            return row.name.charAt(1).toLowerCase()
          }
          return row.name.charAt(0).toLowerCase()
        },
        'Shared on': function (row: IncomingShareResource) {
          localStorage.setItem('grouping-shared-with-me', 'Shared on')
          const recently = Date.now() - 604800000
          const lastMonth = Date.now() - 2592000000
          if (Date.parse(row.sdate) < lastMonth) {
            return 'Older'
          }
          if (Date.parse(row.sdate) >= recently) {
            return 'Recently'
          } else {
            return 'Last month'
          }
        },
        'Share owner': function (row: IncomingShareResource) {
          localStorage.setItem('grouping-shared-with-me', 'Share owner')
          return row?.owner?.displayName
        },
        None: function () {
          localStorage.setItem('grouping-shared-with-me', 'None')
        }
      },
      sortGroups: {
        'Name alphabetically': function (groups: { name: string }[]) {
          // sort in alphabetical order by group name
          const sortedGroups = groups.sort(function (a, b) {
            if (a.name < b.name) {
              return -1
            }
            if (a.name > b.name) {
              return 1
            }
            return 0
          })
          // if sorting is done by name, reverse groups depending on asc/desc
          if (unref(sortBy) === 'name' && unref(sortDir) === 'desc') {
            sortedGroups.reverse()
          }
          return sortedGroups
        },
        'Shared on': function (groups: { name: string }[]) {
          // sort in order: 1-Recently, 2-Last month, 3-Older
          const sortedGroups = []
          const options = ['Recently', 'Last month', 'Older']
          for (const o of options) {
            const found = groups.find((el) => el.name.toLowerCase() === o.toLowerCase())
            if (found) {
              sortedGroups.push(found)
            }
          }
          // if sorting is done by sdate, reverse groups depending on asc/desc
          if (unref(sortBy) === 'sdate' && unref(sortDir) === 'asc') {
            sortedGroups.reverse()
          }
          return sortedGroups
        },
        'Share owner': function (groups: { name: string }[]) {
          // sort in alphabetical order by group name
          const sortedGroups = groups.sort(function (a, b) {
            if (a.name < b.name) {
              return -1
            }
            if (a.name > b.name) {
              return 1
            }
            return 0
          })
          // if sorting is done by owner, reverse groups depending on asc/desc
          if (unref(sortBy) === 'owner' && unref(sortDir) === 'desc') {
            sortedGroups.reverse()
          }
          return sortedGroups
        }
      }
    }
  })

  return {
    groupingSettings
  }
}
