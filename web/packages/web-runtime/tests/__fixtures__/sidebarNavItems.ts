export default [
  {
    name: 'Personal',
    icon: 'folder',
    route: {
      name: 'files-personal',
      path: '/files/list/all'
    },
    active: true
  },
  {
    name: 'Shares',
    icon: 'share-forward',
    route: {
      name: 'files-shared-with-me',
      path: '/files/list/shared-with-me'
    },
    active: false
  },
  {
    name: 'Deleted files',
    icon: 'delete',
    route: {
      name: 'files-trashbin',
      path: '/files/list/trash-bin'
    },
    active: false
  }
]
