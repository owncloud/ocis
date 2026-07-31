export const FolderViewModeConstants = {
  // FIXME: we have a few places where we still match against hardcoded names, get rid of that and this constants
  name: {
    table: 'resource-table',
    condensedTable: 'resource-table-condensed',
    tiles: 'resource-tiles'
  },
  defaultModeName: 'resource-table',
  queryName: 'view-mode',
  tilesSizeDefault: 2,
  tilesSizeMax: 6,
  tilesSizeQueryName: 'tiles-size'
} as const
